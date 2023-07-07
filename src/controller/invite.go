package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/illacloud/illa-supervisor-backend/src/accesscontrol"
	"github.com/illacloud/illa-supervisor-backend/src/model"
)

// InviteByEmail
// email <> first time be invited 				   -> send invite link.
//       <> has been invited and role not changed. -> re-send invite link.
// 	     <> has been invited but role changed 	   -> disable old invite record, generate new invite linke and send invite link.
// 	     <> suspended 							   -> generate new invite linke and send invite link.
//
// @notes:
// 	- can not invite email already used by team member
//
// Execute Phrase:
// - init phrase
//     - get userID
//     - get teamID
//     - get request
//     - validate request
//     - get teamMember
//     - check attribute
//     - get team
//     - check team permission
//     - get user
// - execute invite phrase
//     - check if email and teamID in invite record (email being invite)
//         - exists, resend, end
//         - exists, but role changed,
//             - update invite role
//             - update team_member role
//             - resend invite, end
//     - check if email used by team member (email was used)
//         - true, team_member normal
//             - error feedback, end
//         - false, resend phrase
//             - new invite email
//             - send invite email
//             - set send status
//             - create invite
//             - create pending team_member
//             - return invite link
// - return

func (controller *Controller) InviteMemberByEmail(c *gin.Context) {
	/* Init Phrase */

	// get userID
	userID, errFetchUserID := controller.GetUserIDFromAuth(c)
	teamID, errFetchTeamID := controller.GetMagicIntParamFromRequest(c, PARAM_TEAM_ID)
	if errFetchUserID != nil || errFetchTeamID != nil {
		return
	}

	// get request body
	req := model.NewInviteMemberByEmailRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}
	req.Init()
	userRole := req.ExportUserRole()

	// validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate request body error: "+err.Error())
		return
	}

	// get teamMember (now user role)
	teamMember, errFetchTeamember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errFetchTeamember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "get team member error: "+errFetchTeamember.Error())
		return
	}

	// check attribute
	// validate now user can generate target user role invite link
	// owner & admin -> can invite admin, editor, viewer
	// editor 	     -> can invite editor, viewer
	// viewer 	     -> can invite viewer
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_INVITE)
	if !attrg.CanInvite(userRole) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not generate this type invite link by system permission.")
		return
	}

	// get team
	team, err := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// check if team invite permission was opend
	if controller.DoesTeamInvitePermissionWasClosed(c, attrg, team) {
		return
	}

	// get user
	user, errInFetchUser := controller.Storage.UserStorage.RetrieveByID(userID)
	if errInFetchUser != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "fetch user error: "+errInFetchUser.Error())
		return
	}

	/* Execute Invite Phrase */

	// - check if email and teamID in invite record
	existsInvite, errInFetchInvite := controller.Storage.InviteStorage.RetrieveInviteByTeamIDAndEmail(teamID, req.ExportEmail())

	// already exists
	if errInFetchInvite == nil {

		// set team identifier (the team identifier won't storage in invite record)
		existsInvite.SetTeamIdentifier(team.GetIdentifier())
		existsInvite.SetHosts(req.ExportHosts())

		// share app invite email
		if req.IsShareAppInvite() {
			existsInvite.SetAppID(req.ExportAppID())
		}

		// - exists, resend, end
		if existsInvite.ExportUserRole() == req.ExportUserRole() {
			if controller.SendInvite(c, existsInvite, team, user, req.ExportRedirectPage()) != nil {
				return
			}
			controller.FeedbackInviteByEmail(c, existsInvite)
			return
		}

		// - exists, but role changed, update role, resend
		if existsInvite.ExportUserRole() != req.ExportUserRole() {
			existsTeamMember, errInFetchExistsTeamMember := controller.Storage.TeamMemberStorage.RetrieveTeamMemberByTeamIDAndID(teamID, existsInvite.ExportTeamMemberID())
			if errInFetchExistsTeamMember != nil {
				return
			}
			existsTeamMember.SetUserRole(req.ExportUserRole())
			existsInvite.SetUserRole(req.ExportUserRole())
			if controller.UpdateInvite(c, existsInvite) != nil {
				return
			}
			if controller.UpdateTeamMember(c, existsTeamMember) != nil {
				return
			}
			if controller.SendInvite(c, existsInvite, team, user, req.ExportRedirectPage()) != nil {
				return
			}
			controller.FeedbackInviteByEmail(c, existsInvite)
			return
		}
	}

	// - check if email target user already in teams
	emailUsed, targetTeamMember := controller.EmailAlreadyUsedByTeamMember(c, req.ExportEmail(), teamID)

	// - true, team_member normal, email used, feedback it.
	if emailUsed && targetTeamMember.IsStatusOK() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_EMAIL_ALREADY_USED, "email already used by your team member.")
		return
	}

	// - false, resend phrase
	// - new invite email
	newInviteEmailLink := model.NewInviteEmailLinkByTeamAndRequest(team, req) // team identifier included

	// share app invite email
	if req.IsShareAppInvite() {
		newInviteEmailLink.SetAppID(req.ExportAppID())
	}

	// storage team_member to database
	if controller.StorageTeamMemberByEmailInvite(c, newInviteEmailLink) != nil {
		return
	}

	// - send invite email
	if controller.SendInvite(c, newInviteEmailLink, team, user, req.ExportRedirectPage()) != nil {
		return
	}

	// - set send status
	newInviteEmailLink.SetEmailStatusSuccess()

	// storage invite link to database
	if controller.StorageInvite(c, newInviteEmailLink) != nil {
		return
	}

	controller.FeedbackInviteByEmail(c, newInviteEmailLink)
	return
}

// Execute Phrase:
// - init phrase
//   - get teamID
//   - get userID
//   - get userRole
//   - fetch teamMember record
//   - check attribute
//   - fetch team record
//   - check team permission
//
// - execute invite phrase
//   - fetch invite record
//   - if invite already exists
//   - true, feedback
//   - generate new invite
//   - create invite record
//   - feedback
func (controller *Controller) GenerateInviteLink(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	userRole, errInGetUserRole := controller.GetIntParamFromRequest(c, PARAM_USER_ROLE)
	redirectPage, _ := controller.TestStringParamFromRequest(c, PARAM_REDIRECT_PAGE)
	if errInGetUserID != nil || errInGetUserRole != nil {
		return
	}

	// optional appID
	appIDSetted := false
	appID, errInGetAppID := controller.TestMagicIntParamFromRequest(c, PARAM_APP_ID)
	if errInGetAppID == nil {
		appIDSetted = true
	}

	// get now user role
	teamMember, errFetchTeamember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errFetchTeamember != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "get team member error: "+errFetchTeamember.Error())
		return
	}

	// validate now user can generate target user role invite link
	// owner & admin -> can invite admin, editor, viewer
	// editor 	     -> can invite editor, viewer
	// viewer 	     -> can invite viewer
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_INVITE)
	if !attrg.CanInvite(userRole) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not generate this type invite link by system permission.")
		return
	}

	// get team by id
	team, err := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// check if team invite permission was opend
	if attrg.DoesNowUserAreEditorOrViewer() {
		if !team.DoesEditorOrViewerCanInviteMember() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "your team admin closed the \"editor and viewer can invite members\" permission.")
			return
		}
	}
	// check team invite link config
	tp := team.ExportTeamPermission()
	if !tp.DoesInviteLinkEnabled() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "your team admin closed the \"join by link\" permission.")
		return
	}

	// check if invite link already exists
	invite, err := controller.Storage.InviteStorage.RetrieveAvaliableInviteLinkByTeamIDAndUserRole(teamID, userRole)

	// invite link already generated, just feedback it.
	if err == nil && invite != nil {
		if appIDSetted {
			invite.SetAppID(appID)
			invite.SetTeamIdentifier(team.GetIdentifier())
		}
		controller.FeedbackInviteByLink(c, invite, redirectPage)
		return
	}

	// invite link not exists, generate invite link
	newInviteLink := model.NewInviteLinkByTeamAndUserRole(team, userRole)

	// storage invite link to database
	if controller.StorageInvite(c, newInviteLink) != nil {
		return
	}

	if appIDSetted {
		newInviteLink.SetAppID(appID)
	}

	// return invite link
	controller.FeedbackInviteByLink(c, newInviteLink, redirectPage)
	return
}

// Execute Phrase:
// - init phrase
//     - get teamID
//     - get userID
//     - get userRole
//     - fetch teamMember record
//     - check attribute
//     - fetch team record
//     - check team permission
// - execute invite phrase
//     - fetch invite record
//     - if invite already exists
//         - true, delete invite record
//     - generate new invite
//     - create invite record

func (controller *Controller) RenewInviteLink(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	userRole, errInGetUserRole := controller.GetIntParamFromRequest(c, PARAM_USER_ROLE)
	redirectPage, _ := controller.TestStringParamFromRequest(c, PARAM_REDIRECT_PAGE)
	if errInGetUserID != nil || errInGetUserRole != nil {
		return
	}

	// optional appID
	appIDSetted := false
	appID, errInGetAppID := controller.TestMagicIntParamFromRequest(c, PARAM_APP_ID)
	if errInGetAppID == nil {
		appIDSetted = true
	}

	// get now user role
	teamMember, errFetchTeamember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errFetchTeamember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "get team member error: "+errFetchTeamember.Error())
		return
	}

	// validate now user can generate target user role invite link
	// owner & admin -> can invite admin, editor, viewer
	// editor 	     -> can invite editor, viewer
	// viewer 	     -> can invite viewer
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_INVITE)
	if !attrg.CanInvite(userRole) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not generate this type invite link by system permission.")
		return
	}

	// get team by id
	team, err := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// check if team invite permission was opend (only effect editor and viewer)
	if attrg.DoesNowUserAreEditorOrViewer() {
		if !team.DoesEditorOrViewerCanInviteMember() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "your team admin closed the \"editor and viewer can invite members\" permission.")
			return
		}

	}
	// check team invite link config
	tp := team.ExportTeamPermission()
	if !tp.DoesInviteLinkEnabled() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "your team admin closed the \"join by link\" permission.")
		return
	}

	// check if invite link already exists
	invite, err := controller.Storage.InviteStorage.RetrieveAvaliableInviteLinkByTeamIDAndUserRole(teamID, userRole)

	// invite link already generated, delete it.
	if err == nil && invite != nil {
		if controller.DeleteOldInvite(c, invite) != nil {
			return
		}
	}

	// generate new invite link
	newInviteLink := model.NewInviteLinkByTeamAndUserRole(team, userRole) // team identifier included

	// storage invite link to database
	if controller.StorageInvite(c, newInviteLink) != nil {
		return
	}

	// invite by app
	if appIDSetted {
		newInviteLink.SetAppID(appID)
	}

	// return invite link
	controller.FeedbackInviteByLink(c, newInviteLink, redirectPage)
	return
}

func (controller *Controller) ConfigInviteLink(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// get request body
	req := model.NewConfigInviteLinkRequest()
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

	// validate user
	teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errInRetrieveTeamMember != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
		return
	}

	// validate user role
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_INVITE)
	if !attrg.CanManage(accesscontrol.ACTION_MANAGE_INVITE_LINK) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// get team by id
	team, err := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// update team permission config
	team.ConfigInviteLinkByRequest(req)
	if err := controller.Storage.TeamStorage.UpdateByID(team); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_TEAM, "update team invite permission error: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// Execute Phrase:
// - init Phrase
//     - get UserID
//     - get inviteLinkHash
//     - construct invite
//     - fetch invite record from database
//     - check if now user email does not match invte email
//     - check team permission
// - execute join phrase
//     - if user already in team
//         - true, feedback, end
//     - ok, not exists
//     - insert user in new team
// 	   - if is email invite link
// 	       - delete invite

func (controller *Controller) JoinByLink(c *gin.Context) {
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// get invite link hash
	inviteLinkHash, errFetchInviteLinkHash := controller.GetStringParamFromRequest(c, PARAM_INVITE_LINK_HASH)
	if errFetchInviteLinkHash != nil {
		return
	}
	// construct invite instance
	invite := model.NewInvite()
	errInParseInviteLink := invite.ImportInviteLink(inviteLinkHash)
	if errInParseInviteLink != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_INVITE_LINK_HASH_FAILED, "parse invite link hash failed: "+errInParseInviteLink.Error())
		return
	}

	// get user
	user, errInFetchUser := controller.Storage.UserStorage.RetrieveByID(userID)
	if errInFetchUser != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "fetch user error: "+errInFetchUser.Error())
		return
	}

	// fetch invite record from storage
	inviteRecord, errFetchInvite := controller.Storage.InviteStorage.RetrieveByUID(invite.ExportUID())
	if errFetchInvite != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_INVITE, "fetch invite error: "+errFetchInvite.Error())
		return
	}

	// check if now user email does not match invte email
	if inviteRecord.IsEmailInviteLink() {
		if inviteRecord.ExportEmail() != user.ExportEmail() {
			controller.FeedbackInternalServerError(c, ERROR_FLAG_INVITE_EMAIL_MISMATCH, "invite email mismatch.")
			return
		}
	}

	// get team
	team, err := controller.Storage.TeamStorage.RetrieveByID(inviteRecord.ExportTeamID())
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// check if team closed invite by link permission
	if inviteRecord.IsInviteLink() {
		// check team permission config
		tp := team.ExportTeamPermission()
		if !tp.DoesInviteLinkEnabled() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "your team admin closed the \"join by link\" permission.")
			return
		}
	}

	// check if now user already in team
	userExists, errInCheckUserInTeam := controller.Storage.TeamMemberStorage.DoesTeamIncludedTargetUser(inviteRecord.ExportTeamID(), userID)
	if errInCheckUserInTeam != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_CHECK_TEAM_MEMBER, "check team member error: "+errInCheckUserInTeam.Error())
		return
	}
	if userExists {
		controller.FeedbackBadRequest(c, ERROR_FLAG_USER_ALREADY_JOINED_TEAM, "you already joined this team.")
		return
	}

	// let user join the new team
	teamMember := model.NewTeamMemberByInviteAndUserID(inviteRecord, userID)
	if _, err := controller.Storage.TeamMemberStorage.Create(teamMember); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_CREATE_TEAM_MEMBER, "create team member error: "+err.Error())
		return
	}

	// update invite link status
	if inviteRecord.IsEmailInviteLink() {
		if controller.DeleteOldInvite(c, inviteRecord) != nil {
			return
		}
	}

	// feedback
	controller.FeedbackOK(c, model.NewMyTeamResponse(team, teamMember))
	return

}
