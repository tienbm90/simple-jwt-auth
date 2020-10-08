package models


type UserRole struct {
	User string `json:"user" form:"user" query:"user"`
	Role string `json:"role" form:"role" query:"role"`
}

