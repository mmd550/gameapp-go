package userhandler

import (
	"github.com/labstack/echo/v5"
)

func (handler Handler) AddUserRoutes(e *echo.Echo, prefix string) {
	routeGroup := e.Group(prefix)
	
	routeGroup.POST("/register", handler.Register)
	routeGroup.POST("/login", handler.Login)
	routeGroup.GET("/profile", handler.GetProfile)
}
