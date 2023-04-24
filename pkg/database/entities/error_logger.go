package entities

type ErrorLogger struct {
	Id           int    `json:"id"`
	ErrorMessage string `json:"error_message"`
}
