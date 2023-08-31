package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/illacloud/illa-supervisor-backend/src/accesscontrol"
	"github.com/illacloud/illa-supervisor-backend/src/authenticator"
	"github.com/illacloud/illa-supervisor-backend/src/model"
)

func (controller *Controller) ValidateAccount(c *gin.Context) {
	authorizationToken, errInGetAuthorizationToken := controller.GetStringParamFromHeader(c, PARAM_AUTHORIZATION_TOKEN)
	if errInGetAuthorizationToken != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, authorizationToken)
	if !validated && errInValidate != nil {
		return
	}

	// validate account
	a := controller.Authenticator
	if _, err := a.ManualAuth(authorizationToken); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_ACCOUNT_FAILED, "validate account failed: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) GetTeamPermission(c *gin.Context) {
	authorizationToken, errInGetAuthorizationToken := controller.GetStringParamFromHeader(c, PARAM_AUTHORIZATION_TOKEN)
	teamID := model.TEAM_DEFAULT_ID
	teamIDString, errInGetTeamIDString := controller.GetStringParamFromRequest(c, PARAM_TEAM_ID)
	if errInGetAuthorizationToken != nil || errInGetTeamIDString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, authorizationToken, teamIDString)
	if !validated && errInValidate != nil {
		return
	}

	// fetch team permission info
	team, err := controller.Storage.TeamStorage.RetrieveByID(teamID)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, model.NewGetTeamPermissionResponse(team))
	return
}

func (controller *Controller) CanAccess(c *gin.Context) {
	authorizationToken, errInGetAuthorizationToken := controller.GetStringParamFromHeader(c, PARAM_AUTHORIZATION_TOKEN)
	userID := model.USER_ROLE_ANONYMOUS
	var errInGetUserID error
	if authorizationToken != accesscontrol.ANONYMOUS_AUTH_TOKEN {
		userID, _, errInGetUserID = authenticator.ExtractUserIDFromToken(authorizationToken)
	}
	teamID := model.TEAM_DEFAULT_ID
	unitType, errInGetUnitType := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_TYPE)
	unitID, errInGetUnitID := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_ID)
	attributeID, errInGetAttributeID := controller.GetMagicIntParamFromRequest(c, PARAM_ATTRIBUTE_ID)

	if errInGetAuthorizationToken != nil || errInGetUserID != nil || errInGetUnitType != nil || errInGetUnitID != nil || errInGetAttributeID != nil {
		return
	}

	teamIDString, errInGetTeamIDString := controller.GetStringParamFromRequest(c, PARAM_TEAM_ID)
	unitTypeString, errInGetUnitTypeString := controller.GetStringParamFromRequest(c, PARAM_UNIT_TYPE)
	unitIDString, errInGetUnitIDString := controller.GetStringParamFromRequest(c, PARAM_UNIT_ID)
	attributeIDString, errInGetAttributeIDString := controller.GetStringParamFromRequest(c, PARAM_ATTRIBUTE_ID)

	if errInGetTeamIDString != nil || errInGetUnitTypeString != nil || errInGetUnitIDString != nil || errInGetAttributeIDString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, authorizationToken, teamIDString, unitTypeString, unitIDString, attributeIDString)
	if !validated && errInValidate != nil {
		return
	}

	// validate user
	teamMemberRole := model.USER_ROLE_ANONYMOUS
	if userID != model.USER_ROLE_ANONYMOUS {
		teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
		if errInRetrieveTeamMember != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
			return
		}
		teamMemberRole = teamMember.ExportUserRole()
	}

	// check attribute
	attrg := accesscontrol.NewAttributeGroup(teamMemberRole, unitType)
	attrg.SetUnitID(unitID)
	if !attrg.CanAccess(attributeID) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) CanManage(c *gin.Context) {
	authorizationToken, errInGetAuthorizationToken := controller.GetStringParamFromHeader(c, PARAM_AUTHORIZATION_TOKEN)
	teamID := model.TEAM_DEFAULT_ID
	userID := model.USER_ROLE_ANONYMOUS
	var errInGetUserID error
	if authorizationToken != accesscontrol.ANONYMOUS_AUTH_TOKEN {
		userID, _, errInGetUserID = authenticator.ExtractUserIDFromToken(authorizationToken)
	}
	unitType, errInGetUnitType := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_TYPE)
	unitID, errInGetUnitID := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_ID)
	attributeID, errInGetAttributeID := controller.GetMagicIntParamFromRequest(c, PARAM_ATTRIBUTE_ID)
	if errInGetAuthorizationToken != nil || errInGetUserID != nil || errInGetUnitType != nil || errInGetUnitID != nil || errInGetAttributeID != nil {
		return
	}

	teamIDString, errInGetTeamIDString := controller.GetStringParamFromRequest(c, PARAM_TEAM_ID)
	unitTypeString, errInGetUnitTypeString := controller.GetStringParamFromRequest(c, PARAM_UNIT_TYPE)
	unitIDString, errInGetUnitIDString := controller.GetStringParamFromRequest(c, PARAM_UNIT_ID)
	attributeIDString, errInGetAttributeIDString := controller.GetStringParamFromRequest(c, PARAM_ATTRIBUTE_ID)

	if errInGetTeamIDString != nil || errInGetUnitTypeString != nil || errInGetUnitIDString != nil || errInGetAttributeIDString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, authorizationToken, teamIDString, unitTypeString, unitIDString, attributeIDString)
	if !validated && errInValidate != nil {
		return
	}

	// validate user
	teamMemberRole := model.USER_ROLE_ANONYMOUS
	if userID != model.USER_ROLE_ANONYMOUS {
		teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
		if errInRetrieveTeamMember != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
			return
		}
		teamMemberRole = teamMember.ExportUserRole()
	}

	// check attribute
	attrg := accesscontrol.NewAttributeGroup(teamMemberRole, unitType)
	attrg.SetUnitID(unitID)
	if !attrg.CanManage(attributeID) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) CanManageSpecial(c *gin.Context) {
	authorizationToken, errInGetAuthorizationToken := controller.GetStringParamFromHeader(c, PARAM_AUTHORIZATION_TOKEN)
	teamID := model.TEAM_DEFAULT_ID
	userID := model.USER_ROLE_ANONYMOUS
	var errInGetUserID error
	if authorizationToken != accesscontrol.ANONYMOUS_AUTH_TOKEN {
		userID, _, errInGetUserID = authenticator.ExtractUserIDFromToken(authorizationToken)
	}
	unitType, errInGetUnitType := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_TYPE)
	unitID, errInGetUnitID := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_ID)
	attributeID, errInGetAttributeID := controller.GetMagicIntParamFromRequest(c, PARAM_ATTRIBUTE_ID)
	if errInGetAuthorizationToken != nil || errInGetUserID != nil || errInGetUnitType != nil || errInGetUnitID != nil || errInGetAttributeID != nil {
		return
	}

	teamIDString, errInGetTeamIDString := controller.GetStringParamFromRequest(c, PARAM_TEAM_ID)
	unitTypeString, errInGetUnitTypeString := controller.GetStringParamFromRequest(c, PARAM_UNIT_TYPE)
	unitIDString, errInGetUnitIDString := controller.GetStringParamFromRequest(c, PARAM_UNIT_ID)
	attributeIDString, errInGetAttributeIDString := controller.GetStringParamFromRequest(c, PARAM_ATTRIBUTE_ID)

	if errInGetTeamIDString != nil || errInGetUnitTypeString != nil || errInGetUnitIDString != nil || errInGetAttributeIDString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, authorizationToken, teamIDString, unitTypeString, unitIDString, attributeIDString)
	if !validated && errInValidate != nil {
		return
	}

	// validate user
	teamMemberRole := model.USER_ROLE_ANONYMOUS
	if userID != model.USER_ROLE_ANONYMOUS {
		teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
		if errInRetrieveTeamMember != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
			return
		}
		teamMemberRole = teamMember.ExportUserRole()
	}

	// check attribute
	attrg := accesscontrol.NewAttributeGroup(teamMemberRole, unitType)
	attrg.SetUnitID(unitID)
	if !attrg.CanManageSpecial(attributeID) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) CanModify(c *gin.Context) {
	authorizationToken, errInGetAuthorizationToken := controller.GetStringParamFromHeader(c, PARAM_AUTHORIZATION_TOKEN)
	teamID := model.TEAM_DEFAULT_ID
	userID := model.USER_ROLE_ANONYMOUS
	var errInGetUserID error
	if authorizationToken != accesscontrol.ANONYMOUS_AUTH_TOKEN {
		userID, _, errInGetUserID = authenticator.ExtractUserIDFromToken(authorizationToken)
	}
	unitType, errInGetUnitType := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_TYPE)
	unitID, errInGetUnitID := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_ID)
	attributeID, errInGetAttributeID := controller.GetMagicIntParamFromRequest(c, PARAM_ATTRIBUTE_ID)
	fromID, errInGetFromID := controller.GetMagicIntParamFromRequest(c, PARAM_FROM_ID)
	toID, errInGetToID := controller.GetMagicIntParamFromRequest(c, PARAM_TO_ID)
	if errInGetAuthorizationToken != nil || errInGetUserID != nil || errInGetUnitType != nil || errInGetUnitID != nil || errInGetAttributeID != nil || errInGetFromID != nil || errInGetToID != nil {
		return
	}

	teamIDString, errInGetTeamIDString := controller.GetStringParamFromRequest(c, PARAM_TEAM_ID)
	unitTypeString, errInGetUnitTypeString := controller.GetStringParamFromRequest(c, PARAM_UNIT_TYPE)
	unitIDString, errInGetUnitIDString := controller.GetStringParamFromRequest(c, PARAM_UNIT_ID)
	attributeIDString, errInGetAttributeIDString := controller.GetStringParamFromRequest(c, PARAM_ATTRIBUTE_ID)
	fromIDString, errInGetFromIDString := controller.GetStringParamFromRequest(c, PARAM_FROM_ID)
	toIDString, errInGetToIDString := controller.GetStringParamFromRequest(c, PARAM_TO_ID)

	if errInGetTeamIDString != nil || errInGetUnitTypeString != nil || errInGetUnitIDString != nil || errInGetAttributeIDString != nil || errInGetFromIDString != nil || errInGetToIDString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, authorizationToken, teamIDString, unitTypeString, unitIDString, attributeIDString, fromIDString, toIDString)
	if !validated && errInValidate != nil {
		return
	}

	// validate user
	teamMemberRole := model.USER_ROLE_ANONYMOUS
	if userID != model.USER_ROLE_ANONYMOUS {
		teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
		if errInRetrieveTeamMember != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
			return
		}
		teamMemberRole = teamMember.ExportUserRole()
	}

	// check attribute
	attrg := accesscontrol.NewAttributeGroup(teamMemberRole, unitType)
	attrg.SetUnitID(unitID)
	if !attrg.CanModify(attributeID, fromID, toID) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}

func (controller *Controller) CanDelete(c *gin.Context) {
	authorizationToken, errInGetAuthorizationToken := controller.GetStringParamFromHeader(c, PARAM_AUTHORIZATION_TOKEN)
	teamID := model.TEAM_DEFAULT_ID
	userID := model.USER_ROLE_ANONYMOUS
	var errInGetUserID error
	if authorizationToken != accesscontrol.ANONYMOUS_AUTH_TOKEN {
		userID, _, errInGetUserID = authenticator.ExtractUserIDFromToken(authorizationToken)
	}
	unitType, errInGetUnitType := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_TYPE)
	unitID, errInGetUnitID := controller.GetMagicIntParamFromRequest(c, PARAM_UNIT_ID)
	attributeID, errInGetAttributeID := controller.GetMagicIntParamFromRequest(c, PARAM_ATTRIBUTE_ID)
	if errInGetAuthorizationToken != nil || errInGetUserID != nil || errInGetUnitType != nil || errInGetUnitID != nil || errInGetAttributeID != nil {
		return
	}

	teamIDString, errInGetTeamIDString := controller.GetStringParamFromRequest(c, PARAM_TEAM_ID)
	unitTypeString, errInGetUnitTypeString := controller.GetStringParamFromRequest(c, PARAM_UNIT_TYPE)
	unitIDString, errInGetUnitIDString := controller.GetStringParamFromRequest(c, PARAM_UNIT_ID)
	attributeIDString, errInGetAttributeIDString := controller.GetStringParamFromRequest(c, PARAM_ATTRIBUTE_ID)

	if errInGetTeamIDString != nil || errInGetUnitTypeString != nil || errInGetUnitIDString != nil || errInGetAttributeIDString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, authorizationToken, teamIDString, unitTypeString, unitIDString, attributeIDString)
	if !validated && errInValidate != nil {
		return
	}

	// validate user
	teamMemberRole := model.USER_ROLE_ANONYMOUS
	if userID != model.USER_ROLE_ANONYMOUS {
		teamMember, errInRetrieveTeamMember := controller.Storage.TeamMemberStorage.RetrieveByTeamIDAndUserID(teamID, userID)
		if errInRetrieveTeamMember != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM_MEMBER, "retrieve team member error: "+errInRetrieveTeamMember.Error())
			return
		}
		teamMemberRole = teamMember.ExportUserRole()
	}

	// check attribute
	attrg := accesscontrol.NewAttributeGroup(teamMemberRole, unitType)
	attrg.SetUnitID(unitID)
	if !attrg.CanDelete(attributeID) {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// feedback
	controller.FeedbackOK(c, nil)
	return
}
