package database

import (
	"time"

	"github.com/google/uuid"
)

type Service interface {
	TestConnection() error
	CloseConnection() error
	MakeMigrations() error
}

type UserAccount struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time

	EMail    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	RegisterIP string `gorm:"not null"`
	LoginIP    string
	LastLogin  time.Time
}

type Crosshair struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	RegistrantID uuid.UUID `gorm:"type:uuid;not null" json:"registrant_id"`
	Code         string    `gorm:"not null" json:"code"`
}
