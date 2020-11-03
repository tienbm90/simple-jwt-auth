package models

import "gorm.io/gorm"

type OauthClient struct {
	gorm.Model
	ClientID     string `json:"client_id" query:"client_id"`
	ClientSecret string `json:"client_secret"`
	Domain       string `json:"domain"`
	RedirectURL  string `json:"redirect_url"`
	UserID       string `json:"user_id"`
}
