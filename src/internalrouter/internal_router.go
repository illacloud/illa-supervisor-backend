package internalrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/illacloud/illa-supervisor-backend/src/authenticator"
	"github.com/illacloud/illa-supervisor-backend/src/controller"
)

type Router struct {
	Controller    *controller.Controller
	Authenticator *authenticator.Authenticator
}

func NewRouter(controller *controller.Controller, authenticator *authenticator.Authenticator) *Router {
	return &Router{
		Controller:    controller,
		Authenticator: authenticator,
	}
}

func (r *Router) RegisterRouters(engine *gin.Engine) {
	routerGroup := engine.Group("/api/v1")

	// init user group
	accessControlRouter := routerGroup.Group("/accessControl")
	dataControlRouter := routerGroup.Group("/dataControl")

	// access control routers
	accessControlRouter.GET("/account/validateResult", r.Controller.ValidateAccount)
	accessControlRouter.GET("/teams/:teamID/unitType/:unitType/unitID/:unitID/attribute/canAccess/:attributeID", r.Controller.CanAccess)
	accessControlRouter.GET("/teams/:teamID/unitType/:unitType/unitID/:unitID/attribute/canManage/:attributeID", r.Controller.CanManage)
	accessControlRouter.GET("/teams/:teamID/unitType/:unitType/unitID/:unitID/attribute/canManageSpecial/:attributeID", r.Controller.CanManageSpecial)
	accessControlRouter.GET("/teams/:teamID/unitType/:unitType/unitID/:unitID/attribute/canModify/:attributeID/from/:fromID/to/:toID", r.Controller.CanModify)
	accessControlRouter.GET("/teams/:teamID/unitType/:unitType/unitID/:unitID/attribute/canDelete/:attributeID", r.Controller.CanDelete)

	// data control routers
	dataControlRouter.GET("/users/:targetUserID", r.Controller.GetTargetUserByInternalRequest)
	dataControlRouter.GET("/users/multi/:targetUserIDs", r.Controller.GetTargetUsersByInternalRequest)
	dataControlRouter.GET("/teams/byIdentifier/:teamIdentifier", r.Controller.GetTargetTeamByIdentifier)
}
