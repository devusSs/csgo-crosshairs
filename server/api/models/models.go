package models

import "time"

type RegisterUser struct {
	EMail      string `json:"e_mail"`
	Password   string `json:"password"`
	AdminToken string `json:"admin_token"`
}

type LoginUser struct {
	EMail    string `json:"e_mail"`
	Password string `json:"password"`
}

type ReturnUser struct {
	CreatedAt time.Time `json:"created_at"`
	EMail     string    `json:"e_mail"`
	Role      string    `json:"role"`
}
