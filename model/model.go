package model

import (
	"any-days.com/celebs/db"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID uint `gorm:"primarykey" json:"id" json:"id,omitempty"`

	Name              string `gorm:"type:varchar(255)" json:"name,omitempty"`
	ApiKey            string `gorm:"type:varchar(255);unique_index"`
	Provider          string `json:"provider"`
	Email             string `json:"email"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	NickName          string `json:"nick_name"`
	Description       string `json:"description"`
	RemoteUserID      string `gorm:"index,unique" json:"-"`
	AvatarURL         string `json:"avatar_url"`
	Location          string `json:"-"`
	AccessToken       string `json:"-"`
	AccessTokenSecret string `json:"-"`
	RefreshToken      string `json:"-"`
	IDToken           string `json:"-"`

	ExpiresAt     time.Time      `json:"-"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	ActivatedAt   *time.Time     `json:"-"`
	DeactivatedAt sql.NullTime   `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func Migrate() {
	err := db.Db().AutoMigrate(
		User{},
	)

	if err != nil {
		panic(fmt.Sprintf("Error while migrating: %s", err))
	}
}
