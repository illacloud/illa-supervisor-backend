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
		controller.FeedbackBadRequest(c, ERROR_FLAG_SEND_VERIFICATION_CODE_FAILED, "send verification code error: "+err.Error())
		return
	}

	// cache token
	errInCacheJWTToken := controller.Cache.JWTCache.SetTokenForEmail(req.Email, vToken)
	if errInCacheJWTToken != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAHCE_JWT_TOKEN_FAILED, "cache jwt token failed error: "+errInCacheJWTToken.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
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
	controller.FeedbackOK(c, model.NewSignInResponse(user))
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

	// fetch verification token from cache
	var errInFetchJWTToken error
	req.VerificationToken, errInFetchJWTToken = controller.Cache.JWTCache.GetTokenByEmail(req.Email)
	if errInFetchJWTToken != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_FETCH_JWT_TOKEN_FROM_CACHE, "fetch jwt token from cache failed: "+errInFetchJWTToken.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_GENERATE_PASSWORD_FAILED, "generate password error: "+err.Error())
		return
	}
	user.SetPasswordByByte(hashPwd)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update user password error: "+err.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CREATE_UPLOAD_URL_FAILED, "get upload URL failed: "+errInGetPreSignedURL.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user Nickname
	user.SetNickname(req.Nickname)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update user error: "+err.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user Nickname
	user.SetAvatar(req.Avatar)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update user error: "+err.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_GENERATE_PASSWORD_FAILED, "generate password error: "+err.Error())
		return
	}
	user.SetPasswordByByte(hashPwd)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update password error: "+err.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user language
	user.SetLanguage(req.Language)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update language error: "+err.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// update user language
	user.SetIsTutorialViewed(req.IsTutorialViewed)
	if err := controller.Storage.UserStorage.UpdateByID(user); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_USER, "update isTutorialViewed error: "+err.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, nil)
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
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
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_DELETE_USER, "delete user by id error: "+errInDeleteUser.Error())
		return
	}

	// delete from team member
	errInDeleteTeamMember := controller.Storage.TeamMemberStorage.DeleteByUserID(userID)
	if errInDeleteTeamMember != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_DELETE_TEAM_MEMBER, "delete user from team member by user id error: "+errInDeleteTeamMember.Error())
		return
	}

	// ok, feedback
	controller.FeedbackOK(c, nil)
	return
}
