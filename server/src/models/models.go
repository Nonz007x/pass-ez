package models

import "gorm.io/gorm"

type (
	Fact struct {
		gorm.Model
		Question string `json:"question" gorm:"text;not null;default:null"`
		Answer   string `json:"answer" gorm:"text;not null;default:null"`
	}

	User struct {
		Id       string `gorm:"not null;unique"`
		Email    string `gorm:"not null;unique"`
		MasterPassword string `gorm:"not null;"`
	}

	RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ErrorResponse struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_descripton"`
		Message          string `json:"message"`
	}
)
