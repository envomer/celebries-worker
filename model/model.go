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

// TMDB Movie
type Movie struct {
	ID uint `gorm:"primarykey" json:"id" json:"id,omitempty"`
}

type Genre struct {
	ID     uint   `gorm:"primarykey" json:"id" json:"id,omitempty"`
	TmdbID int    `json:"tmdb_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

type MovieGenre struct {
	MovieID uint `gorm:"index" json:"movie_id,omitempty"`
	GenreID uint `gorm:"index" json:"genre_id,omitempty"`
}

type MovieCast struct {
	ID        uint   `gorm:"primarykey" json:"id" json:"id,omitempty"`
	MovieID   uint   `gorm:"index" json:"movie_id,omitempty"`
	Character string `json:"character,omitempty"`
	Order     int    `json:"order,omitempty"`
	// Person    *Person `gorm:"foreignKey:PersonID" json:"person,omitempty"`
	PersonID uint `gorm:"index" json:"person_id,omitempty"`
}

type Person struct {
	ID                 uint    `gorm:"primarykey" json:"id" json:"id,omitempty"`
	TmdbID             int     `gorm:"index,unique" json:"tmdb_id,omitempty"`
	Name               string  `json:"name,omitempty"`
	Birthday           string  `json:"birthday,omitempty"`
	Deathday           string  `json:"deathday,omitempty"`
	PlaceOfBirth       string  `json:"place_of_birth,omitempty"`
	ProfilePath        string  `json:"profile_path,omitempty"`
	Adult              bool    `json:"adult,omitempty"`
	KnownForDepartment string  `json:"known_for_department,omitempty"`
	Popularity         float32 `json:"popularity"`
}

func Migrate() {
	err := db.Db().AutoMigrate(
		//User{},
		//Movie{},
		//Genre{},
		//MovieGenre{},
		//MovieCast{},
		Person{},
	)

	if err != nil {
		panic(fmt.Sprintf("Error while migrating: %s", err))
	}
}
