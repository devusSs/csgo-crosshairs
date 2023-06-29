package database

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetPostgresVersion() (string, error)

	TestConnection() error
	CloseConnection() error
	MakeMigrations() error

	CreateNewEngineerToken(string) error
	GetLatestEngineerToken() (string, error)

	AddUser(*UserAccount) (*UserAccount, error)
	GetUserByVerificationCode(*UserAccount) (*UserAccount, error)
	UpdateUserVerification(*UserAccount) (*UserAccount, error)
	GetUserByEmail(*UserAccount) (*UserAccount, error)
	UpdateUserLogin(*UserAccount) (*UserAccount, error)
	GetUserByUID(*UserAccount) (*UserAccount, error)
	UpdateUserCrosshairCount(*UserAccount) (*UserAccount, error)
	AddResetPasswordCodeAndTime(*UserAccount) (*UserAccount, error)
	GetUserByResetpasswordCode(*UserAccount) (*UserAccount, error)
	UpdateUserPassword(*UserAccount) (*UserAccount, error)
	UpdateUserPasswordRaw(*UserAccount) (*UserAccount, error)
	UpdateVerifyMailResendTime(*UserAccount) (*UserAccount, error)
	UpdateUserAvatarURL(*UserAccount) (*UserAccount, error)

	AddUserTwitchDetails(*UserAccount) (*UserAccount, error)
	GetUserByTwitchLogin(*UserAccount) (*UserAccount, error)

	AddCrosshair(*Crosshair) (*Crosshair, error)
	GetAllCrosshairsFromUser(uuid.UUID) ([]*Crosshair, error)
	GetAllCrosshairsFromUserSortByDate(uuid.UUID) ([]*Crosshair, error)
	DeleteAllCrosshairsFromUser(uuid.UUID) error
	DeleteCrosshairFromUserByCode(uuid.UUID, string) error
	EditCrosshairNote(*Crosshair) (*Crosshair, error)

	GetAllUsers() ([]*UserAccount, error)
	GetAllCrosshairs() ([]*Crosshair, error)

	AddEvent(*Event) (*Event, error)
	GetEvents() ([]*Event, error)
	GetEventsByType(string) ([]*Event, error)
	GetEventsWithLimit(int) ([]*Event, error)
	GetEventsByTypeWithLimit(string, int) ([]*Event, error)

	WriteTwitchBotLog(*TwitchBotLog) error
	GetAllTwitchBotLogEntries() ([]*TwitchBotLog, error)
	GetLatestTwitchBotLogWithLimit(int) ([]*TwitchBotLog, error)
	GetLatestTwitchBotLogByType(string) ([]*TwitchBotLog, error)
	GetLatestTwitchBotLogByTypeWithLimit(string, int) ([]*TwitchBotLog, error)

	AddTwitchTokenRefreshStore(*TwitchRefreshTokenStore) (*TwitchRefreshTokenStore, error)
	GetLatestTwitchTokenRefreshStore(*TwitchRefreshTokenStore) (*TwitchRefreshTokenStore, error)
	DeleteAllTwitchTokenRefreshStore(*TwitchRefreshTokenStore) error
}

type EngineerToken struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time

	Token string `gorm:"not null"`
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

	RequestNewVerifyMailTime time.Time

	PasswordResetCode     string
	PasswordResetCodeTime time.Time

	AvatarURL string 

	RegisterIP string `gorm:"not null"`
	LoginIP    string
	LastLogin  time.Time

	TwitchID        string 
	TwitchLogin     string 
	TwitchCreatedAt time.Time

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
	CreatedAt time.Time
	Type      EventType `gorm:"embedded;not null" json:"type"`
	Data      EventData `gorm:"embedded;not null" json:"data"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
}

// Submodels for Event struct.
type EventType string

const (
	UserRegistered      EventType = "user_registered"
	UserChangedPassword EventType = "user_password_change"
	UserUploadedAvatar  EventType = "user_uploaded_avatar"
)

type EventData struct {
	URL      string `json:"url"`
	Method   string `json:"method"`
	IssuerIP string `json:"issuer"`
}

type TwitchBotLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"` // Should be a JSON string.
	Issuer    string    `json:"issuer"`
}

func MarshalTwitchBotLogMessage(message interface{}) (string, error) {
	bytes, err := json.Marshal(message)
	return string(bytes), err
}

type TwitchBotLogIssuer string

const (
	Root             TwitchBotLogIssuer = "root"
	JoinChannel      TwitchBotLogIssuer = "join_channel"
	HandleStatus     TwitchBotLogIssuer = "handle_status"
	HandleCrosshairs TwitchBotLogIssuer = "handle_crosshairs"
)

type TwitchRefreshTokenStore struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time

	TwitchID             string
	TwitchLogin          string
	RefreshToken         string
	RefreshTokenAcquired time.Time
	AccessToken          string
	AccessTokenExpiry    time.Time
}
