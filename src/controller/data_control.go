package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/illacloud/illa-supervisior-backend/src/model"
)

func (controller *Controller) GetTargetUserByInternalRequest(c *gin.Context) {
	targetUserID, errInGetTargetUserID := controller.GetMagicIntParamFromRequest(c, PARAM_TARGET_USER_ID)
	targetUserIDString, errInGetTargetUserIDString := controller.GetStringParamFromRequest(c, PARAM_TARGET_USER_ID)
	if errInGetTargetUserID != nil || errInGetTargetUserIDString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, targetUserIDString)
	if !validated && errInValidate != nil {
		return
	}

	// fetch target user info
	user, err := controller.Storage.UserStorage.RetrieveByID(targetUserID)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, model.NewGetTargetUserByInternalRequestResponse(user))
	return
}

func (controller *Controller) GetTargetTeamByIdentifier(c *gin.Context) {
	teamIdentifier, errInGetTeamIdentifier := controller.GetStringParamFromRequest(c, PARAM_TARGET_TEAM_IDENTIFIER)
	if errInGetTeamIdentifier != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, teamIdentifier)
	if !validated && errInValidate != nil {
		return
	}

	// fetch target team info
	team, err := controller.Storage.TeamStorage.RetrieveByIdentifier(teamIdentifier)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get team error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, model.NewGetTargetTeamByInternalRequestResponse(team))
	return
}
