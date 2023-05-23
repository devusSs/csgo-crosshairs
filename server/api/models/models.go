package models

import (
	"time"

	"github.com/google/uuid"
)

// Request models
type RegisterUser struct {
	EMail    string `json:"e_mail"`
	Password string `json:"password"`
}

type LoginUser struct {
	EMail    string `json:"e_mail"`
	Password string `json:"password"`
}

type AddCrosshair struct {
	Code string `json:"code"`
	Note string `json:"note"`
}

type ResetPassword struct {
	EMail string `json:"e_mail"`
}

type ResetPasswordFinal struct {
	Password string `json:"password"`
}

// Response models
type ReturnUser struct {
	CreatedAt time.Time `json:"created_at"`
	EMail     string    `json:"e_mail"`
	Role      string    `json:"role"`
}

type ReturnUserAdmin struct {
	ID                   uuid.UUID `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	EMail                string    `json:"e_mail"`
	Role                 string    `json:"role"`
	VerifiedMail         bool      `json:"verified_mail"`
	RegisterIP           string    `json:"register_ip"`
	LoginIP              string    `json:"login_ip"`
	LastLogin            time.Time `json:"last_login"`
	CrosshairsRegistered int       `json:"crosshairs_registered"`
}

type MultipleUsersAdmin struct {
	Users []ReturnUserAdmin `json:"users"`
}

type Crosshair struct {
	ID    uuid.UUID `json:"id"`
	Added time.Time `json:"added"`
	Code  string    `json:"code"`
	Note  string    `json:"note"`
}

type GetMultipleCrosshairs struct {
	Crosshairs []Crosshair `json:"crosshairs"`
}

type GetOneCrosshair struct {
	Crosshair Crosshair `json:"crosshair"`
}
