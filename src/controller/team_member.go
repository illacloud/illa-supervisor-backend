package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/illacloud/illa-supervisor-backend/src/accesscontrol"
	"github.com/illacloud/illa-supervisor-backend/src/model"
)

func (controller *Controller) GetAllTeamMember(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// validate user
	teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errInRetrieveTeamMember != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you are not in this team.")
		return
	}

	// validate user role
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_TEAM_MEMBER)
	if !attrg.CanAccess(accesscontrol.ACTION_ACCESS_VIEW) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// get team by id
	teamMembers, err := controller.Storage.TeamMemberStorage.RetrieveByTeamID(teamID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "get team members error: "+err.Error())
		return
	}

	// get team members detail info
	userIDs := model.PickUpUserIDsInUserMembers(teamMembers)
	users, err := controller.Storage.UserStorage.RetrieveByIDs(userIDs)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get team members detail error: "+err.Error())
		return
	}

	// get pending user emails
	teamMemberIDs := model.PickUpTeamMemberIDsInTeamMembers(teamMembers)
	invites, _ := controller.Storage.InviteStorage.RetrieveInviteByEmailByIDs(teamMemberIDs)

	userForExportLT := model.BuildLookUpTableForUserExport(users)
	inviteForExportLT := model.BuildLookUpTableForInvitesExport(invites)
	teamMembersForExport, err := controller.assembleTeamMembers(teamMembers, userForExportLT, inviteForExportLT)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_BUILD_TEAM_MEMBER_LIST_FAILED, "build team members list error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, model.NewGetAllTeamMembersResponse(teamMembersForExport))
	return
}

func (controller *Controller) assembleTeamMembers(teamMembers []*model.TeamMember, userExportLT map[int]*model.UserForExport, inviteForExportLT map[int]*model.InviteForExport) ([]*model.TeamMemberWithUserInfoForExport, error) {
	// sort
	sort.Slice(teamMembers, func(i, j int) bool {
		return teamMembers[i].CreatedAt.After(teamMembers[j].CreatedAt)
	})
	// assemble
	var r []*model.TeamMemberWithUserInfoForExport
	for _, teamMember := range teamMembers {
		targetUser, found := userExportLT[teamMember.UserID]
		if !found {
			// is pending user ?
			if teamMember.UserID == model.PENDING_USER_ID {
				// fetch email from invite
				inviteForExport, found := inviteForExportLT[teamMember.ID]
				if !found {
					fmt.Printf("[error] can not found a team member from invite record.")
					continue
				}
				// new pending user
				pendingUser, errInNewPendingUser := model.NewPendingUserByInviteForExport(inviteForExport)
				exportedPendingUser := pendingUser.Export()
				exportedPendingUser.SetTeamMemberID(teamMember.ID)
				if errInNewPendingUser != nil {
					return nil, errInNewPendingUser
					continue
				}
				// append
				r = append(r, teamMember.ExportWithUserInfo(exportedPendingUser))
				continue
			}

			// user can not found
			return nil, errors.New("found a team member \"" + strconv.Itoa(teamMember.UserID) + "\" but can not found user detail info.")
		}
		targetUser.SetTeamMemberID(teamMember.ID)
		r = append(r, teamMember.ExportWithUserInfo(targetUser))
	}
	return r, nil
}

func (controller *Controller) GetTeamMember(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	targetUserID, errInGetTargetUserID := controller.GetMagicIntParamFromRequest(c, PARAM_TARGET_USER_ID)
	if errInGetUserID != nil || errInGetTargetUserID != nil {
		return
	}

	// validate user
	teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errInRetrieveTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
		return
	}

	// validate user role
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_TEAM_MEMBER)
	if !attrg.CanAccess(accesscontrol.ACTION_ACCESS_VIEW) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// get target team member by id
	targetTeamMember, errInRetrieveTargetTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, targetUserID)
	if errInRetrieveTargetTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTargetTeamMember.Error())
		return
	}

	// get team members detail info
	user, err := controller.Storage.UserStorage.RetrieveByID(targetUserID)
	if err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_USER, "get team members detail error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, model.NewGetTeamMemberResponse(targetTeamMember.ExportWithUserInfo(user.Export())))
	return
}

// global role config
// owner  <MODIFY> owner                 -> any.                   NO, owner must be transferred before modify himself (a team must have an owner).
// 		  <MODIFY> admin, editor, viewer -> owner.                 YES, it's transfer owner action (when the owner transferred, the old owner became an admin).
//        <MODIFY> admin, editor, viewer -> admin, editor, viewer. YES.
// admin  <MODIFY> owner                 -> any.                   NO.
// 		  <MODIFY> admin, editor, viewer -> owner.                 NO.
//        <MODIFY> admin, editor, viewer -> admin, editor, viewer. YES.
// editor <MODIFY> any                   -> any.                   NO.
// viewer <MODIFY> any                   -> any.                   NO.
func (controller *Controller) UpdateTeamMemberRole(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	targetTeamMemberID, errInGetTargetTeamMemberID := controller.GetMagicIntParamFromRequest(c, PARAM_TARGET_TEAM_MEMBER_ID)
	if errInGetUserID != nil || errInGetTargetTeamMemberID != nil {
		return
	}

	// get request body
	req := model.NewUpdateTeamMemberRoleRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate payload required fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// get team member
	teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errInRetrieveTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
		return
	}

	// check if editor and viewer can manage team member in this team
	team, errInRetrieveTeam := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if errInRetrieveTeam != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "retrieve team error: "+errInRetrieveTeam.Error())
		return
	}
	tp := team.ExportTeamPermission()
	if teamMember.IsEditor() {
		if !tp.DoesEditorCanManageTeamMember() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "you can not manage team member due to team permission settings.")
			return
		}
	}
	if teamMember.IsViewer() {
		if !tp.DoesViewerCanManageTeamMember() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_TEAM_CLOSED_THE_PERMISSION, "you can not manage team member due to team permission settings.")
			return
		}
	}

	// get target team member
	targetTeamMember, errInRetrieveTargetTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndID(teamID, targetTeamMemberID)
	if errInRetrieveTargetTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "get target team member error: "+errInRetrieveTargetTeamMember.Error())
		return
	}

	// validate user role
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_TEAM_MEMBER)
	if !attrg.CanManage(accesscontrol.ACTION_MANAGE_ROLE) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// check if now user can manage target user's role to target role
	if !attrg.CanModifyRoleFromTo(targetTeamMember.ExportUserRole(), req.ExportUserRole()) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// manage user himself role action, by default team owner can not set himself to non-onwer role
	targetUserID := targetTeamMember.ExportUserID()
	if targetUserID == userID {
		if teamMember.IsOwner() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_OWNER_ROLE_MUST_BE_TRANSFERED, "you must transfer your team owner role before modify your role.")
			return
		}
	}

	// fetch invite record
	existsInvite, errInFetchInvite := controller.Storage.InviteStorage.RetrieveAvaliableInviteByTeamIDAndTeamMemberID(teamID, targetTeamMember.ExportID())

	// check if owner transfer owner to a pending user (invite record exists)
	if errInFetchInvite == nil {
		if req.IsTransferOwner() {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_TRANSFER_OWNER_TO_PENDING_USER, "can not transfet team owner to a pending member.")
			return
		}
	}

	// transfer owner action, check if now user is owner and target user is owner himself.
	if req.IsTransferOwner() {
		// set owner to admin
		teamMember.UpdateTeamMemberRole(accesscontrol.USER_ROLE_ADMIN)
	}

	// update target teammeber
	targetTeamMember.UpdateByUpdateTeamMemberRoleRequest(req)
	if err := controller.Storage.TeamMemberStorage.Update(targetTeamMember); err != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_TEAM_MEMBER, "update target team member error: "+err.Error())
		return
	}

	// if is pending user (invite record exists), update invite
	if errInFetchInvite == nil {
		existsInvite.SetUserRole(req.ExportUserRole())
		if controller.UpdateInvite(c, existsInvite) == nil {
			return
		}
	}

	// update owner when needed
	if req.IsTransferOwner() {
		if err := controller.Storage.TeamMemberStorage.Update(teamMember); err != nil {
			controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_UPDATE_TEAM_MEMBER, "update team member error: "+err.Error())
			return
		}
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}

// remove team member config
// owner  <DELETE> owner                 NO, owner can not delete himself (a team must have an owner).
// 		  <DELETE> admin, editor, viewer YES.
// admin  <DELETE> owner                 NO.
// 		  <DELETE> admin, editor, viewer YES.
// editor <DELETE> any                   NO.
// 		  <DELETE> editor, viewer        YES.
// viewer <DELETE> any                   NO.
// 		  <DELETE> viewer 				 YES.

func (controller *Controller) RemoveTeamMember(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	targetTeamMemberID, errInGetTargetTeamMemberID := controller.GetMagicIntParamFromRequest(c, PARAM_TARGET_TEAM_MEMBER_ID)
	if errInGetUserID != nil || errInGetTargetTeamMemberID != nil {
		return
	}

	// validate user
	teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errInRetrieveTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
		return
	}

	// get target team member
	targetTeamMember, errInRetrieveTargetTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndID(teamID, targetTeamMemberID)
	if errInRetrieveTargetTeamMember != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "get target team member error: "+errInRetrieveTargetTeamMember.Error())
		return
	}

	// validate user role
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_TEAM_MEMBER)
	if !attrg.CanManage(accesscontrol.ACTION_MANAGE_REMOVE_MEMBER) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// check if removing owner
	if targetTeamMember.IsOwner() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_REMOVE_OWNER_FROM_TEAM, "can not remove an owner from a team.")
		return
	}

	// remove user from the team
	errInDelete := controller.Storage.TeamMemberStorage.DeleteByIDAndTeamID(targetTeamMemberID, teamID)
	if errInDelete != nil {
		controller.FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_DELETE_TEAM_MEMBER, "delete team member by team id and user id error: "+errInDelete.Error())
		return
	}

	// remove user from invite for pending user
	controller.Storage.InviteStorage.DeleteByTeamIDAndTeamMemberID(teamID, targetTeamMember.ExportID())

	// remove target user
	controller.Storage.UserStorage.DeleteByID(targetTeamMember.ExportUserID())

	// feedback
	controller.FeedbackOK(c, nil)
	return
}
