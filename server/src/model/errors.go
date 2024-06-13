package model

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_descripton"`
	Message          string `json:"message"`
}

var (
	InternalServerError = ErrorResponse{
		Error:            "500",
		ErrorDescription: "internal_server_error",
		Message:          "Something went wrong. Please try again later.",
	}

	JsonParsingError = ErrorResponse{
		Error:            "400",
		ErrorDescription: "cannot_parse_json",
		Message:          "The request body contains invalid JSON. Please try again.",
	}

	EmailConflictError = ErrorResponse{
		Error:            "409",
		ErrorDescription: "email_already_in_use",
		Message:          "Email is already in use. Try again with a different email.",
	}

	UserNotFoundError = ErrorResponse{
		Error:            "404",
		ErrorDescription: "user_not_found",
		Message:          "Wrong email or password. Try again.",
	}

	InvalidTokenError = ErrorResponse{
		Error:            "401",
		ErrorDescription: "invalid_token",
		Message:          "Unauthorized",
	}

	TokenExpiredError = ErrorResponse{
		Error:            "401",
		ErrorDescription: "token_expired",
		Message:          "Unauthorized",
	}
)
