package controller

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/illacloud/illa-supervisor-backend/src/model"
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

func (controller *Controller) GetTargetUsersByInternalRequest(c *gin.Context) {
	targetUserIDsString, errInGetTargetUserIDsString := controller.GetStringParamFromRequest(c, PARAM_TARGET_USER_IDS)
	if errInGetTargetUserIDsString != nil {
		return
	}

	// validate request data
	validated, errInValidate := controller.ValidateRequestTokenFromHeader(c, targetUserIDsString)
	if !validated && errInValidate != nil {
		return
	}

	// convert ids from string to int
	targetUserIDsStringSlice := strings.Split(targetUserIDsString, ",")
	var targetUserIDsInInt []int
	for _, idInString := range targetUserIDsStringSlice {
		idInInt, _ := strconv.Atoi(idInString)
		targetUserIDsInInt = append(targetUserIDsInInt, idInInt)
	}

	// fetch target user info
	users, err := controller.Storage.UserStorage.RetrieveByIDs(targetUserIDsInInt)
	if err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_USER, "get user error: "+err.Error())
		return
	}

	// feedback
	controller.FeedbackOK(c, model.NewGetTargetUsersByInternalRequestResponse(users))
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
