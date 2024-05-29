package models

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_descripton"`
	Message          string `json:"message"`
}

var DatabaseError = ErrorResponse{
	Error:            "internal_server_error",
	ErrorDescription: "database_error",
	Message:          "something went wrong. Try again.",
}