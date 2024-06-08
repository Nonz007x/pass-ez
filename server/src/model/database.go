package model

type (
	User struct {
		Id             string `gorm:"not null;unique"`
		Email          string `gorm:"not null;unique"`
		MasterPassword string `gorm:"not null;"`
		Salt           string `gorm:"not null;"`
	}

	UserVault struct {
		VaultId string `gorm:"not null;"`
		UserId  string `gorm:"not null;"`
	}

	Vault struct {
		Id      string `gorm:"not null;unique"`
		Key     string `gorm:"not null;"`
		OwnerId string `gorm:"not null;"`
	}

	SaltKey struct {
		Salt string
		Key  string
	}
)

func (UserVault) TableName() string {
	return "user_vault"
}
