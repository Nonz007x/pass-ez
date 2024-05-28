package models

type (

	User struct {
		Id             string `gorm:"not null;unique"`
		Email          string `gorm:"not null;unique"`
		MasterPassword string `gorm:"not null;"`
		Salt           string `gorm:"not null;"`
	}

	RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Salt     string `json:"salt"`
		VaultKey string `json:"vault_key"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Vault struct {
		Id      string `gorm:"not null;unique"`
		Key     string `gorm:"not null;"`
		OwnerId string `gorm:"not null;"`
	}

	SaltKey struct {
		Salt string `json:"salt"`
		Key string `json:"protected_key"`
	}

	UserVault struct {
		VaultId string `gorm:"not null;"`
		UserId  string `gorm:"not null;"`
	}

	SaltResponse struct {
		Email   string `json:"email"`
		Salt string `json:"salt"`
	}

	ErrorResponse struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_descripton"`
		Message          string `json:"message"`
	}
)

func (UserVault) TableName() string {
	return "user_vault"
}

var DatabaseError = ErrorResponse{
	Error:            "internal_server_error",
	ErrorDescription: "database_error",
	Message:          "something went wrong. Try again.",
}