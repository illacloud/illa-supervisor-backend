package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/illacloud/illa-supervisor-backend/src/accesscontrol"
	"github.com/illacloud/illa-supervisor-backend/src/model"
)

func (controller *Controller) GetMyTeams(c *gin.Context) {
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// retrieve
	teamMembers, errInGetTeamMember := controller.Storage.TeamMemberStorage.RetrieveByUserID(userID)
	if errInGetTeamMember != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team by id error: "+errInGetTeamMember.Error())
		return
	}
	teamIDs := model.PickUpTeamIDsInTeamMembers(teamMembers)
	teams, errInGetTeam := controller.Storage.TeamStorage.RetrieveByIDs(teamIDs)
	if errInGetTeam != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "retrieve team by ids error: "+errInGetTeam.Error())
		return
	}

	// build lookup table for feedback
	tmlt := model.BuildTeamIDLookUpTableForTeamMemberExport(teamMembers)

	// feedback
	controller.FeedbackOK(c, model.NewGetMyTeamsResponse(teams, tmlt))
	return
}

func (controller *Controller) UpdateTeamConfig(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// get request body
	var rawRequest map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&rawRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate user
	teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errInRetrieveTeamMember != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "please make sure that your can access this team. retrieve team member error: "+errInRetrieveTeamMember.Error())
		return
	}

	// validate user role
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_TEAM)
	if !attrg.CanManage(accesscontrol.ACTION_MANAGE_TEAM_CONFIG) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// get team by id
	team, err := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// update team config
	errInConstructRawConfig := team.UpdateByUpdateTeamConfigRawRequest(rawRequest)
	if errInConstructRawConfig != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_BUILD_TEAM_CONFIG_FAILED, "build team config error: "+errInConstructRawConfig.Error())
		return
	}

	// update
	if err := controller.Storage.TeamStorage.UpdateByID(team); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_TEAM, "update team error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) UpdateTeamPermission(c *gin.Context) {
	// get team id & user id
	teamID := model.TEAM_DEFAULT_ID
	userID, errInGetUserID := controller.GetUserIDFromAuth(c)
	if errInGetUserID != nil {
		return
	}

	// get request body
	var rawRequest map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&rawRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error: "+err.Error())
		return
	}

	// validate user
	teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
	if errInRetrieveTeamMember != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "please make sure that your can access this team. retrieve team member error: "+errInRetrieveTeamMember.Error())
		return
	}

	// validate user role
	attrg := accesscontrol.NewAttributeGroup(teamMember.ExportUserRole(), accesscontrol.UNIT_TYPE_TEAM)
	if !attrg.CanManage(accesscontrol.ACTION_SPECIAL_EDITOR_AND_VIEWER_CAN_INVITE_BY_LINK_SW) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// get team by id
	team, err := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// update team permission
	errInParseRawReq := team.UpdateByUpdateTeamPermissionRawRequest(rawRequest)
	if errInParseRawReq != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_BUILD_TEAM_PERMISSION_FAILED, "build team permission error: "+errInParseRawReq.Error())
		return
	}
	if err := controller.Storage.TeamStorage.UpdateByID(team); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_UPDATE_TEAM, "update team error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}
