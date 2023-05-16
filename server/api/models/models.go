package models

type RegisterUser struct {
	EMail      string `json:"e_mail"`
	Password   string `json:"password"`
	AdminToken string `json:"admin_token"`
}

type LoginUser struct {
	EMail    string `json:"e_mail"`
	Password string `json:"password"`
}
