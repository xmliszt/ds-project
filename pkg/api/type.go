package api

import "github.com/labstack/echo/v4"

type Response struct {
	Success bool
	Error   string
	Data    interface{}
}

type Route struct {
	Method  string
	Path    string
	Handler echo.HandlerFunc
	Auth    bool
}
