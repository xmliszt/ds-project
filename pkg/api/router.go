package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouterBuilder struct{}

func (rb *RouterBuilder) New() *echo.Echo {
	return echo.New()
}

func GetRouter() *echo.Echo {
	routerBuilder := &RouterBuilder{}
	router := routerBuilder.New()

	apiRouter := router.Group("/api/v1")
	dashboardRouter := router.Group("/dashboard")

	// Logger middleware
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} From: ${remote_ip} Method: ${method} URI: ${uri} Status: ${status} Error: ${error} Latency: ${latency_human}\n",
	}))

	// Bind JWT token auth to sub-routers
	apiRouter.Use(middleware.JWT([]byte("secret")))
	dashboardRouter.Use(middleware.JWT([]byte("secret")))

	// Views
	router.GET("/", Home)
	dashboardRouter.GET("", Dashboard)

	// User
	router.POST("/register", Register)
	router.POST("/login", LogIn)

	// Secret
	apiRouter.GET("/secret", GetSecret)
	apiRouter.PUT("/secret", PutSecret)
	apiRouter.DELETE("/secret", DeleteSecret)
	apiRouter.GET("/secrets", GetAllSecrets)

	return router
}
