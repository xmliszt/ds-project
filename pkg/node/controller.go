package node

import (
	"github.com/labstack/echo/v4"
	"github.com/xmliszt/e-safe/pkg/api"
)

func (n *Node) getRoutes() []api.Route {
	return []api.Route{
		{
			Method:  echo.GET,
			Path:    "/secret",
			Handler: n.getSecret,
			Auth:    true,
		},
		{
			Method:  echo.PUT,
			Path:    "/secret",
			Handler: n.putSecret,
			Auth:    true,
		},
		{
			Method:  echo.DELETE,
			Path:    "/secret",
			Handler: n.deleteSecret,
			Auth:    true,
		},
		{
			Method:  echo.GET,
			Path:    "/secrets",
			Handler: n.getAllSecrets,
			Auth:    true,
		},
		{
			Method:  echo.POST,
			Path:    "/register",
			Handler: n.register,
			Auth:    false,
		},
		{
			Method:  echo.POST,
			Path:    "/login",
			Handler: n.logIn,
			Auth:    false,
		},
		{
			Method:  echo.GET,
			Path:    "/monitor",
			Handler: n.getMonitorInfo,
			Auth:    true,
		},
	}
}
