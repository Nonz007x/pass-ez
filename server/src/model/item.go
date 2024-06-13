package model

import (
	"time"
)

type (
	Item struct {
		ID      string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
		VaultID string    `gorm:"type:varchar(255);not null" json:"vault_id"`
		Name    string    `gorm:"type:varchar(1000);not null" json:"name"`
		Note    string    `gorm:"type:varchar(10000)" json:"note"`
		Type    int16     `gorm:"type:smallint;not null" json:"type"`
		Created time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created"`
		Updated time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated"`
	}

	Login struct {
		ItemID               string            `gorm:"primaryKey;type:varchar(255)" json:"item_id"`
		Username             string            `gorm:"type:varchar(1000)" json:"username"`
		Password             string            `gorm:"type:varchar(1000)" json:"password"`
		PasswordRevisionDate time.Time         `gorm:"type:timestamptz" json:"password_revision_date"`
		PasswordHistories    []PasswordHistory `gorm:"foreignKey:LoginID" json:"password_histories"`
		URIs                 []URI             `gorm:"foreignKey:LoginID" json:"uris"`
	}

	PasswordHistory struct {
		ID       int64     `gorm:"primaryKey;autoIncrement" json:"id"`
		LoginID  string    `gorm:"type:varchar(255);not null" json:"login_id"`
		Password string    `gorm:"type:varchar(10000)" json:"password"`
		LastUsed time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"last_used"`
	}

	URI struct {
		ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
		LoginID string `gorm:"type:varchar(255);not null" json:"login_id"`
		URI     string `gorm:"type:varchar(10000)" json:"uri"`
	}
)
