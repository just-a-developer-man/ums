package router

import (
	"ums/internal/api/v1/middleware"

	"github.com/gin-gonic/gin"
)

const (
	apiGroup   = "/api"
	v1Group    = "/v1"
	user       = "/user"
	adminGroup = "/admin"
	password   = "/password"
	email      = "/email"
	name       = "/name"
	users      = "/users"
	id         = "/:id"
)

type Handler interface {
	CreateUser(c *gin.Context)
	CreateUserWithRole(c *gin.Context)
	ReadUser(c *gin.Context)
	UpdateUserData(c *gin.Context)
	UpdateUserPassword(c *gin.Context)
	UpdateUserEmail(c *gin.Context)
	UpdateUserName(c *gin.Context)
	DeleteUser(c *gin.Context)
	ReadUsers(c *gin.Context)
}

func SetupRouter(h Handler, mode string) *gin.Engine {
	gin.SetMode(mode)

	router := gin.New()
	router.Use(middleware.Logging())
	setupApiGroup(router, h)

	return router
}

func setupApiGroup(r *gin.Engine, h Handler) {
	api := r.Group(apiGroup)
	{
		setupV1Subgroup(api, h)
	}
}

func setupV1Subgroup(g *gin.RouterGroup, h Handler) {
	v1 := g.Group(v1Group)
	{
		setupUserSubgroup(v1, h)
		setupAdminSubgroup(v1, h)
	}
}

func setupAdminSubgroup(g *gin.RouterGroup, h Handler) {
	admin := g.Group(adminGroup)
	admin.Use(middleware.AdminAuth())
	{
		admin.GET(users, h.ReadUsers)
		admin.POST(user, h.CreateUserWithRole)
		setupIdSubgroup(admin, h)
	}
}

func setupUserSubgroup(g *gin.RouterGroup, h Handler) {
	user := g.Group(user)
	{
		user.POST("", h.CreateUser)
		setupIdSubgroup(user, h)
	}
}

func setupIdSubgroup(g *gin.RouterGroup, h Handler) {
	idG := g.Group(id)
	{
		idG.GET("", h.ReadUser)
		idG.PUT("", h.UpdateUserData)
		idG.PATCH(password, h.UpdateUserPassword)
		idG.PATCH(email, h.UpdateUserEmail)
		idG.PATCH(name, h.UpdateUserName)
		idG.DELETE("", h.DeleteUser)
	}
}
