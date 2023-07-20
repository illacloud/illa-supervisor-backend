package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/illacloud/illa-supervisor-backend/src/authenticator"
	"github.com/illacloud/illa-supervisor-backend/src/model"
)

// generate and send verification code email
func (controller *Controller) GetVerificationCode(c *gin.Context) {
	// resolve request
	req := model.NewVerificationRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate request body error: "+err.Error())
		return
	}

	vToken, err := model.GenerateAndSendVerificationCode(req.Email, req.Usage)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_SEND_VERIFICATION_CODE_FAILED, "send verification code error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, model.NewGetVerificationCodeResponse(vToken))
	return
}

// user sign-up
// signup  							    -> check if email was used
//
//	<> with email invite token    -> check invite_token
//	                              -  create user
//	                              -  update team_member
//	                              -  delete invite_token
//	<> with link invite token     -> check invite_token
//	                              -  create team_member
//	                              -  create user
//	<> without token 				-> create user
func (controller *Controller) SignUp(c *gin.Context) {
	// get request body
	req := model.NewSignUpRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate request body error: "+err.Error())
		return
	}

	// check if team setting `blockRegister` is true
	team, errInGetTeam := controller.Storage.TeamStorage.RetrieveByID(DEFAULT_TEAM_ID)
	if errInGetTeam != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team failed: "+errInGetTeam.Error())
		return
	}
	teamPermission := team.ExportTeamPermission()
	if teamPermission.DoesBlockRegister() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_REGISTER_BLOCKED, "signup failed: registration blocked")
		return
	}

	// check if email laready used
	_, errFetchUserRecord := controller.Storage.UserStorage.RetrieveByEmail(req.ExportEmail())
	if errFetchUserRecord == nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_EMAIL_HAS_BEEN_TAKEN, "email already has been taken.")
		return
	}

	// @todo: validate verification code when user smtp server configureable

	// [join team if invite token setted]
	if req.IsSignUpWithInviteLink() {
		// construct invite token
		invite := model.NewInvite()
		errInParseInviteLink := invite.ImportInviteLink(req.ExportInviteToken())
		if errInParseInviteLink != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_INVITE_LINK_HASH_FAILED, "parse invite link hash failed: "+errInParseInviteLink.Error())
			return
		}

		// fetch invite record from storage
		inviteRecord, errFetchInvite := controller.Storage.InviteStorage.RetrieveByUID(invite.ExportUID())
		if errFetchInvite != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_INVITE, "fetch invite record error: "+errFetchInvite.Error())
			return
		}

		// check if invite link are avaliable (unused & status are avaliable)
		if !inviteRecord.IsAvaliable() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_INVITATION_LINK_UNAVALIABLE, "invite link unavaliable.")
			return
		}

		// [is invite by link]
		if inviteRecord.IsInviteLink() {
			controller.signUpWithLinkToken(req, inviteRecord, c)
			return
		}

		// [else if, invite by email]
		if inviteRecord.IsEmailInviteLink() {
			controller.signUpWithEmailToken(req, inviteRecord, c)
			return
		}
	}
	// [else, create new user]
	controller.signUpWithoutToken(req, c)
	return
}

// - execute invite phrase
//   - check request email match invite record email
//   - create user
//   - fetch teamMember record
//   - set teamMember active
//   - update teamMember record
//   - delete invite record
//   - generate access token
//   - feedback
func (controller *Controller) signUpWithEmailToken(req *model.SignUpRequest, inviteRecord *model.Invite, c *gin.Context) {
	// check if signup request email does not match invite record email
	if req.ExportEmail() != inviteRecord.ExportEmail() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_UP_EMAIL_MISMATCH, "signup email does not match invite email.")
		return
	}

	// ok, create new user
	user, err := model.NewUserBySignUpRequest(req)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_BUILD_USER_INFO_FAILED, "user customized info construct error: "+err.Error())
		return
	}
	newUserIDInt, errInCreateUser := controller.Storage.UserStorage.Create(user)
	if errInCreateUser != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_CREATE_USER, "create user error: "+errInCreateUser.Error())
		return
	}
	if !req.IsSubscribed {
		_ = model.SendSubscriptionEmail(req.Email)
	}
	user.SetID(newUserIDInt)

	// update team member
	pendingTeamMember, errInFetchPendingTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndID(inviteRecord.ExportTeamID(), inviteRecord.ExportTeamMemberID())
	if errInFetchPendingTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "fetch pending team member record error: "+errInFetchPendingTeamMember.Error())
		return
	}
	pendingTeamMember.SetUserID(user.ExportID())
	pendingTeamMember.ActiveUser()
	if err := controller.Storage.TeamMemberStorage.Update(pendingTeamMember); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_TEAM_MEMBER, "update pending team member info error: "+err.Error())
		return
	}

	// delete invite link
	if controller.DeleteOldInvite(c, inviteRecord) != nil {
		return
	}

	// generate access token and refresh token
	accessToken, _ := model.CreateAccessToken(user.ID, user.UID)
	expiredAtString, errInExtract := authenticator.ExtractExpiresAtFromTokenInString(accessToken)
	if errInExtract != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "check token expired at failed")
		return
	}
	errInCacheTokenExpiredAt := controller.Cache.JWTCache.InitUserJWTTokenExpiredAt(user, expiredAtString)
	if errInCacheTokenExpiredAt != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "cache token expired at failed")
		return
	}
	c.Header("illa-token", accessToken)

	// ok, feedback
	controller.FeedbackOK(c, model.NewSignUpResponse(user))
	return
}

func (controller *Controller) signUpWithLinkToken(req *model.SignUpRequest, inviteRecord *model.Invite, c *gin.Context) {
	// check if team closed invite by link permission
	// get team by id
	team, err := controller.Storage.TeamStorage.RetrieveByID(inviteRecord.ExportTeamID())
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// check team permission config
	tp := team.ExportTeamPermission()
	if !tp.DoesInviteLinkEnabled() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "your team admin closed the \"join by link\" permission.")
		return
	}

	// ok, create user
	user, err := model.NewUserBySignUpRequest(req)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_BUILD_USER_INFO_FAILED, "user customized info construct error: "+err.Error())
		return
	}
	newUserIDInt, errInCreateUser := controller.Storage.UserStorage.Create(user)
	if errInCreateUser != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_CREATE_USER, "create user error: "+errInCreateUser.Error())
		return
	}
	if !req.IsSubscribed {
		_ = model.SendSubscriptionEmail(req.Email)
	}
	user.SetID(newUserIDInt)

	// let user join the new team
	teamMember := model.NewTeamMemberByInviteAndUserID(inviteRecord, newUserIDInt)
	if _, err := controller.Storage.TeamMemberStorage.Create(teamMember); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_CREATE_TEAM_MEMBER, "create team member error: "+err.Error())
		return
	}

	// generate access token and refresh token
	accessToken, _ := model.CreateAccessToken(user.ID, user.UID)
	expiredAtString, errInExtract := authenticator.ExtractExpiresAtFromTokenInString(accessToken)
	if errInExtract != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "check token expired at failed")
		return
	}
	errInCacheTokenExpiredAt := controller.Cache.JWTCache.InitUserJWTTokenExpiredAt(user, expiredAtString)
	if errInCacheTokenExpiredAt != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "cache token expired at failed")
		return
	}
	c.Header("illa-token", accessToken)

	// ok, feedback
	controller.FeedbackOK(c, model.NewSignUpResponse(user))
	return
}

func (controller *Controller) signUpWithoutToken(req *model.SignUpRequest, c *gin.Context) {
	// eliminate duplicate user
	if duplicateUser, _ := controller.Storage.UserStorage.RetrieveByEmail(req.Email); duplicateUser != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_EMAIL_HAS_BEEN_TAKEN, "duplicate email address")
		return
	}

	// construct user
	user, err := model.NewUserBySignUpRequest(req)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_BUILD_USER_INFO_FAILED, "user customized info construct error: "+err.Error())
		return
	}

	// create user
	newUserIDInt, errInCreateUser := controller.Storage.UserStorage.Create(user)
	if errInCreateUser != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_USER, "create user error: "+errInCreateUser.Error())
		return
	}
	if !req.IsSubscribed {
		_ = model.SendSubscriptionEmail(req.Email)
	}
	user.SetID(newUserIDInt)

	// create team member
	teamMember := model.NewEditorTeamMemberByUserID(newUserIDInt)
	if _, err := controller.Storage.TeamMemberStorage.Create(teamMember); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_CREATE_TEAM_MEMBER, "create team member error: "+err.Error())
		return
	}

	// generate access token and refresh token
	accessToken, _ := model.CreateAccessToken(user.ID, user.UID)
	expiredAtString, errInExtract := authenticator.ExtractExpiresAtFromTokenInString(accessToken)
	if errInExtract != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "check token expired at failed")
		return
	}
	errInCacheTokenExpiredAt := controller.Cache.JWTCache.InitUserJWTTokenExpiredAt(user, expiredAtString)
	if errInCacheTokenExpiredAt != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "cache token expired at failed")
		return
	}
	c.Header("illa-token", accessToken)

	// ok, feedback
	controller.FeedbackOK(c, model.NewSignUpResponse(user))
	return
}

// user sign-in
func (controller *Controller) SignIn(c *gin.Context) {
	// get request body
	req := model.NewSignInRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate request body error: "+err.Error())
		return
	}

	// fetch user by email
	user, err := controller.Storage.UserStorage.RetrieveByEmail(req.Email)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "invalid email or password")
		return
	}

	// validate password with password digest
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(req.Password))
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "invalid email or password")
		return
	}

	// generate access token and refresh token
	accessToken, _ := model.CreateAccessToken(user.ID, user.UID)
	expiredAtString, errInExtract := authenticator.ExtractExpiresAtFromTokenInString(accessToken)
	if errInExtract != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "check token expired at failed")
		return
	}
	errInCacheTokenExpiredAt := controller.Cache.JWTCache.InitUserJWTTokenExpiredAt(user, expiredAtString)
	if errInCacheTokenExpiredAt != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "cache token expired at failed")
		return
	}
	c.Header("illa-token", accessToken)

	// ok, feedback
	controller.FeedbackOK(c, model.NewSignUpResponse(user))
	return
}

func (controller *Controller) ForgetPassword(c *gin.Context) {
	// get request body
	req := model.NewForgetPasswordRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// fetch user by email
	user, err := controller.Storage.UserStorage.RetrieveByEmail(req.Email)
	if err != nil || user.ID == 0 {
		controller.FeedbackBadRequest(c, ERROR_FLAG_NO_SUCH_USER, "no such user")
		return
	}

	// validate verification code
	validCode, err := model.ValidateVerificationCode(req.VerificationCode, req.VerificationToken,
		req.Email, "forgetpwd")
	if !validCode || err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_VERIFICATION_CODE_FAILED, "validate verification code error: "+err.Error())
		return
	}

	// update password
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_GENERATE_PASSWORD_FAILED, "generate password error: "+err.Error())
		return
	}
	user.SetPasswordByByte(hashPwd)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update user password error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) Logout(c *gin.Context) {
	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}
	// clean jwt expireAt token
	errInDeleteTokenExpiredAt := controller.Cache.JWTCache.CleanUserJWTTokenExpiredAt(user)
	if errInDeleteTokenExpiredAt != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_SIGN_IN_FAILED, "clean token expired at cache failed")
		return
	}
	// @todo: logout method
	c.JSON(http.StatusOK, nil)
	return
}

func (controller *Controller) GetUserAvatarUploadAddress(c *gin.Context) {
	fileName, errInGetFineName := controller.GetStringParamFromRequest(c, PARAM_FILE_NAME)
	if errInGetFineName != nil {
		return
	}

	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// get upload address
	systemDrive := model.NewSystemDrive(controller.Drive)
	systemDrive.SetUser(user)
	resignedURL, errInGetPreSignedURL := systemDrive.GetUserAvatarUploadPreSignedURL(fileName)
	if errInGetPreSignedURL != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CREATE_UPLOAD_URL_FAILED, "get upload URL failed: "+errInGetPreSignedURL.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, model.NewUserAvatarUploadAddressResponse(resignedURL))
	return
}

func (controller *Controller) UpdateNickname(c *gin.Context) {
	// get request body
	req := model.NewUpdateNicknameRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user Nickname
	user.SetNickname(req.Nickname)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update user error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, model.NewUpdateUserResponse(user))
	return
}

func (controller *Controller) UpdateAvatar(c *gin.Context) {
	// get request body
	req := model.NewUpdateAvatarRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user Nickname
	user.SetAvatar(req.Avatar)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update user error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, model.NewUpdateUserResponse(user))
	return
}

func (controller *Controller) UpdatePassword(c *gin.Context) {
	// get request body
	req := model.NewChangePasswordRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// validate current password with password digest
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(req.CurrentPassword)); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PASSWORD_INVALIED, "current password incorrect")
		return
	}

	// update password
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_GENERATE_PASSWORD_FAILED, "generate password error: "+err.Error())
		return
	}
	user.SetPasswordByByte(hashPwd)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update password error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) UpdateLanguage(c *gin.Context) {
	// get request body
	req := model.NewUpdateLanguageRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate request body error: "+err.Error())
		return
	}

	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user language
	user.SetLanguage(req.Language)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update language error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, model.NewUpdateUserResponse(user))
	return
}

func (controller *Controller) UpdateIsTutorialViewed(c *gin.Context) {
	// get request body
	req := model.NewUpdateIsTutorialViewedRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate request body error: "+err.Error())
		return
	}

	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user language
	user.SetIsTutorialViewed(req.IsTutorialViewed)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update isTutorialViewed error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) CreateUser(c *gin.Context) {
	// get request body
	CreateUserRequest := model.NewCreateUserRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&CreateUserRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate CreateUserRequest required fields
	validate := validator.New()
	if err := validate.Struct(CreateUserRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// create user
	User := model.NewUserByCreateUserRequest(CreateUserRequest)
	userID, err := controller.Storage.UserStorage.Create(User)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_CREATE_USER, "create user error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, model.NewCreateUserResponse(userID))
	return
}

func (controller *Controller) RetrieveUserByID(c *gin.Context) {
	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// retrieve
	user, err := controller.Storage.UserStorage.RetrieveByID(userID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, model.NewGetUserByIDResponse(user))
	return
}

func (controller *Controller) DeleteUser(c *gin.Context) {
	// get user by id
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// check if user is owner
	nowUserIsTeamOwner, errInCheckIsOwner := controller.Storage.TeamMemberStorage.IsNowUserIsTeamOwner(userID)
	if errInCheckIsOwner != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "check toeam owner failed: "+errInCheckIsOwner.Error())
		return
	}
	if nowUserIsTeamOwner {
		controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_MUST_TRANSFERED_BEFORE_USER_SUSPEND, "please transfer your team owner role, then you can suspend your account.")
		return
	}

	// delete
	errInDeleteUser := controller.Storage.UserStorage.DeleteByID(userID)
	if errInDeleteUser != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_DELETE_USER, "delete user by id error: "+errInDeleteUser.Error())
		return
	}

	// delete from team member
	errInDeleteTeamMember := controller.Storage.TeamMemberStorage.DeleteByUserID(userID)
	if errInDeleteTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_DELETE_TEAM_MEMBER, "delete user from team member by user id error: "+errInDeleteTeamMember.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, nil)
	return
}
