package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HomePage
func Home(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "This is Home Page!")
}

// HomePage
func Dashboard(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "This is user's Dashboard!")
}
