package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/illacloud/illa-supervisior-backend/src/accesscontrol"
	"github.com/illacloud/illa-supervisior-backend/src/model"
)

func (controller *Controller) EmailAlreadyUsedByTeamMember(c *gin.Context, email string, teamID int) (bool, *model.TeamMember) {
	user, errInFetchUser := controller.Storage.UserStorage.RetrieveByEmail(email)
	if errInFetchUser == nil {
		teamMember, errFetchTeamember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, user.ExportID())
		if errFetchTeamember == nil {
			return true, teamMember
		}
	}
	return false, nil
}

func (controller *Controller) DoesTeamInvitePermissionWasClosed(c *gin.Context, attrg *accesscontrol.AttributeGroup, team *model.Team) bool {
	if attrg.DoesNowUserAreEditorOrViewer() {
		if !team.DoesEditorOrViewerCanInviteMember() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "your team admin closed the \"editor and viewer can invite members\" permission.")
			return true
		}
	}
	return false
}

func (controller *Controller) SendInvite(c *gin.Context, invite *model.Invite, team *model.Team, user *model.User) error {
	// email share app invite
	if invite.IsShareAppInvite() {
		m := model.NewEmailShareAppMessage(invite, team, user)
		token := controller.RequestTokenValidator.GenerateValidateToken(m.UserName, m.TeamName, m.TeamIcon, m.Email, m.AppLink, m.Language)
		m.SetValidateToken(token)
		errInSendEmail := model.SendShareAppEmail(m)
		if errInSendEmail != nil {
			invite.SetEmailStatusFailed()
			controller.FeedbackInternalServerError(c, ERROR_FLAG_SEND_EMAIL_FAILED, "send invite email failed."+errInSendEmail.Error())
			return errInSendEmail
		}
		return nil
	}

	// email invite
	m := model.NewEmailInviteMessage(invite, team, user)
	token := controller.RequestTokenValidator.GenerateValidateToken(m.UserName, m.TeamName, m.TeamIcon, m.Email, m.JoinLink, m.Language)
	m.SetValidateToken(token)
	errInSendEmail := model.SendInviteEmail(m)
	if errInSendEmail != nil {
		invite.SetEmailStatusFailed()
		controller.FeedbackInternalServerError(c, ERROR_FLAG_SEND_EMAIL_FAILED, "send invite email failed."+errInSendEmail.Error())
		return errInSendEmail
	}
	return nil
}

func (controller *Controller) StorageInvite(c *gin.Context, invite *model.Invite) error {
	_, errInCreateInvite := controller.Storage.InviteStorage.Create(invite)
	if errInCreateInvite != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CREATE_LINK_FAILED, "create invite link error: "+errInCreateInvite.Error())
		return errInCreateInvite
	}
	return nil
}

func (controller *Controller) UpdateInvite(c *gin.Context, invite *model.Invite) error {
	invite.InitUpdatedAt()
	errInUpdateInvite := controller.Storage.InviteStorage.Update(invite)
	if errInUpdateInvite != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_INVITE, "update exists invite record failed: "+errInUpdateInvite.Error())
		return errInUpdateInvite
	}
	return nil
}

func (controller *Controller) UpdateTeamMember(c *gin.Context, teamMember *model.TeamMember) error {
	teamMember.InitUpdatedAt()
	errInUpdateTeamMember := controller.Storage.TeamMemberStorage.Update(teamMember)
	if errInUpdateTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_TEAM_MEMBER, "update team member record failed: "+errInUpdateTeamMember.Error())
		return errInUpdateTeamMember
	}
	return nil
}

func (controller *Controller) DeleteOldInvite(c *gin.Context, invite *model.Invite) error {
	errInDeleteInvite := controller.Storage.InviteStorage.DeleteByID(invite.ExportID())
	if errInDeleteInvite != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_DELETE_INVITE, "remove old invite link error: "+errInDeleteInvite.Error())
		return errInDeleteInvite
	}
	return nil
}

func (controller *Controller) StorageTeamMemberByEmailInvite(c *gin.Context, invite *model.Invite) error {
	newTeamMember := model.NewPendingTeamMemberByInvite(invite)
	teamMemberID, errInCreateTeamMember := controller.Storage.TeamMemberStorage.Create(newTeamMember)
	if errInCreateTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_CREATE_INVITE, "create invite link error: "+errInCreateTeamMember.Error())
		return errInCreateTeamMember
	}
	invite.SetTeamMemberID(teamMemberID)
	return nil
}

func (controller *Controller) FeedbackNewInviteByEmail(c *gin.Context, invite *model.Invite) {
	resp := model.NewInviteMemberByEmailResponseByInviteRecord(invite, "brandnew invite link generated and send by email.")
	controller.FeedbackOK(c, resp)
}

func (controller *Controller) FeedbackInviteByEmail(c *gin.Context, invite *model.Invite) {
	resp := model.NewInviteMemberByEmailResponseByInviteRecord(invite, "resend invite email.")
	controller.FeedbackOK(c, resp)
}

func (controller *Controller) FeedbackInviteByLink(c *gin.Context, invite *model.Invite) {
	resp := model.NewGenerateInviteLinkResponseByInvite(invite)
	controller.FeedbackOK(c, resp)
}
