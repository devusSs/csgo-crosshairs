package database

import (
	"time"

	"github.com/google/uuid"
)

type Service interface {
	TestConnection() error
	CloseConnection() error
	MakeMigrations() error

	AddUser(*UserAccount) (*UserAccount, error)
	GetUserByVerificationCode(*UserAccount) (*UserAccount, error)
	UpdateUserVerification(*UserAccount) (*UserAccount, error)
	GetUserByEmail(*UserAccount) (*UserAccount, error)
	UpdateUserLogin(*UserAccount) (*UserAccount, error)
	GetUserByUID(*UserAccount) (*UserAccount, error)
}

type UserAccount struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time

	EMail            string `gorm:"unique;not null"`
	Password         string `gorm:"not null"`
	Role             string `gorm:"not null"`
	VerificationCode string `gorm:"not null"`
	VerifiedMail     bool

	RegisterIP string `gorm:"not null"`
	LoginIP    string
	LastLogin  time.Time
}

type Crosshair struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	RegistrantID uuid.UUID `gorm:"type:uuid;not null" json:"registrant_id"`
	Code         string    `gorm:"not null" json:"code"`
}
