package models


type UserRole struct {
	User string `json:"user" forms:"user" query:"user"`
	Role string `json:"role" forms:"role" query:"role"`
}

