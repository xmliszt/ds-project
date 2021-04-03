package node

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xmliszt/e-safe/pkg/api"
)

type RouterBuilder struct{}

func (rb *RouterBuilder) New() *echo.Echo {
	return echo.New()
}

func (n *Node) getRouter() *echo.Echo {
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

	routes := n.getRoutes()

	for _, route := range routes {
		switch route.Method {
		case echo.GET:
			if route.Auth {
				apiRouter.GET(route.Path, route.Handler)
			} else {
				router.GET(route.Path, route.Handler)
			}
		case echo.POST:
			if route.Auth {
				apiRouter.POST(route.Path, route.Handler)
			} else {
				router.POST(route.Path, route.Handler)
			}
		case echo.PUT:
			if route.Auth {
				apiRouter.PUT(route.Path, route.Handler)
			} else {
				router.PUT(route.Path, route.Handler)
			}
		case echo.DELETE:
			if route.Auth {
				apiRouter.DELETE(route.Path, route.Handler)
			} else {
				router.DELETE(route.Path, route.Handler)
			}
		}

	}

	// Views
	router.GET("/", api.Home)
	dashboardRouter.GET("", api.Dashboard)

	return router
}
