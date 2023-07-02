package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"transac_api/lib/structs"
	"transac_api/service/user"
)

type Router struct {
	engine      *gin.Engine
	userService *user.Service
}

func MustMakeRouter(engine *gin.Engine, userService *user.Service) *Router {
	return &Router{
		engine:      engine,
		userService: userService,
	}
}

func (router Router) AddRoutes() {
	group := router.engine.Group("")
	{
		group.GET("/users", router.GetUsers)
	}
}

func (router Router) GetUsers(c *gin.Context) {
	u, err := router.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.GetInternalServerError())
		return
	}

	c.JSON(http.StatusOK, u)
}
