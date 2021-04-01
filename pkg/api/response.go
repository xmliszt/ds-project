package api

type Response struct {
	Success bool
	Error   string
	Data    interface{}
}
