package models

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Salt     string `json:"salt"`
		VaultKey string `json:"vault_key"`
	}

	Token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
