package router

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

	authRouter := routerGroup.Group("/auth")
	usersRouter := routerGroup.Group("/users")
	teamsRouter := routerGroup.Group("/teams")
	joinRouter := routerGroup.Group("/join")
	statusRouter := routerGroup.Group("/status")

	// register auth
	usersRouter.Use(r.Authenticator.JWTAuth())
	teamsRouter.Use(r.Authenticator.JWTAuth())
	joinRouter.Use(r.Authenticator.JWTAuth())

	// auth routers
	authRouter.POST("/verification", r.Controller.GetVerificationCode)
	authRouter.POST("/signup", r.Controller.SignUp)
	authRouter.POST("/signin", r.Controller.SignIn)
	authRouter.POST("/forgetPassword", r.Controller.ForgetPassword)

	// user routers
	usersRouter.GET("", r.Controller.RetrieveUserByID)
	usersRouter.GET("/avatar/uploadAddress/fileName/:fileName", r.Controller.GetUserAvatarUploadAddress)
	usersRouter.PATCH("/password", r.Controller.UpdatePassword)
	usersRouter.PATCH("/nickname", r.Controller.UpdateNickname)
	usersRouter.PATCH("/avatar", r.Controller.UpdateAvatar)
	usersRouter.PATCH("/language", r.Controller.UpdateLanguage)
	usersRouter.PATCH("/tutorialViewed", r.Controller.UpdateIsTutorialViewed)
	usersRouter.DELETE("", r.Controller.DeleteUser)
	usersRouter.POST("/logout", r.Controller.Logout)

	// teams routers
	teamsRouter.GET("/my", r.Controller.GetMyTeams)
	teamsRouter.PATCH("/:teamID/config", r.Controller.UpdateTeamConfig)
	teamsRouter.PATCH("/:teamID/permission", r.Controller.UpdateTeamPermission)

	// team members
	teamsRouter.GET("/:teamID/members", r.Controller.GetAllTeamMember)
	teamsRouter.GET("/:teamID/users/:targetUserID", r.Controller.GetTeamMember)
	teamsRouter.PATCH("/:teamID/teamMembers/:targetTeamMemberID/role", r.Controller.UpdateTeamMemberRole)
	teamsRouter.DELETE("/:teamID/teamMembers/:targetTeamMemberID", r.Controller.RemoveTeamMember)

	// invite routers
	teamsRouter.PATCH("/:teamID/configInviteLink", r.Controller.ConfigInviteLink)
	teamsRouter.GET("/:teamID/inviteLink/userRole/:userRole", r.Controller.GenerateInviteLink)
	teamsRouter.GET("/:teamID/newInviteLink/userRole/:userRole", r.Controller.RenewInviteLink)
	teamsRouter.POST("/:teamID/inviteByEmail", r.Controller.InviteMemberByEmail)

	// join routers
	joinRouter.PUT("/:inviteLinkHash", r.Controller.JoinByLink)

	// share routers
	teamsRouter.GET("/:teamID/shareAppLink/userRole/:userRole/apps/:appID", r.Controller.GenerateInviteLink)
	teamsRouter.GET("/:teamID/shareAppLink/userRole/:userRole/apps/:appID/redirectPage/:redirectPage", r.Controller.GenerateInviteLink)
	teamsRouter.GET("/:teamID/newShareAppLink/userRole/:userRole/apps/:appID", r.Controller.RenewInviteLink)
	teamsRouter.GET("/:teamID/newShareAppLink/userRole/:userRole/apps/:appID/redirectPage/:redirectPage", r.Controller.RenewInviteLink)
	teamsRouter.POST("/:teamID/shareAppByEmail", r.Controller.InviteMemberByEmail)

	// status router
	statusRouter.GET("", r.Controller.Status)

}
