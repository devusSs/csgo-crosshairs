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
	UpdateUserCrosshairCount(*UserAccount) (*UserAccount, error)
	AddResetPasswordCode(*UserAccount) (*UserAccount, error)
	GetUserByResetpasswordCode(*UserAccount) (*UserAccount, error)
	UpdateUserPassword(*UserAccount) (*UserAccount, error)

	AddCrosshair(*Crosshair) (*Crosshair, error)
	GetAllCrosshairsFromUser(uuid.UUID) ([]*Crosshair, error)
	DeleteAllCrosshairsFromUser(uuid.UUID) error
	DeleteCrosshairFromUserByCode(uuid.UUID, string) error
	EditCrosshairNote(*Crosshair) (*Crosshair, error)

	GetAllUsers() ([]*UserAccount, error)
	GetAllCrosshairs() ([]*Crosshair, error)

	AddEvent(*Event) (*Event, error)
	GetEvents() ([]*Event, error)
	GetEventsByType(string) ([]*Event, error)
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

	PasswordResetCode string

	RegisterIP string `gorm:"not null"`
	LoginIP    string
	LastLogin  time.Time

	// For now we will only allow 20 crosshairs per user.
	CrosshairsRegistered int
}

type Crosshair struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time

	RegistrantID uuid.UUID `gorm:"type:uuid;not null"`
	Code         string    `gorm:"not null"`
	Note         string

	RegisterIP string `gorm:"not null"`
}

type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type      EventType `gorm:"embedded;not null" json:"type"`
	Data      EventData `gorm:"embedded;not null" json:"data"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
}

// Submodels for Event struct.
type EventType string

// TODO: add more events?
const (
	UserRegistered      EventType = "user_registered"
	UserChangedPassword EventType = "user_password_change"
)

type EventData struct {
	URL      string `json:"url"`
	Method   string `json:"method"`
	IssuerIP string `json:"issuer"`
}
