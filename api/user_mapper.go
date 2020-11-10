package api

import "github.com/simple-jwt-auth/models"

type UserDTO struct {
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Email    string `gorm:"size:100;not null;unique" json:"email"`
	Password string `gorm:"size:100;not null;" json:"password"`
}

func ToUser(userDTO UserDTO) models.User {
	return models.User{UserName: userDTO.Username, Email: userDTO.Email, Password: userDTO.Password}
}

func ToUserDTO(user models.User) UserDTO {
	return UserDTO{Username: user.UserName, Email: user.Email, Password: user.Password}
}

func ToUserDTOs(users []models.User) []UserDTO {
	userDTOS := make([]UserDTO, len(users))

	for i, itm := range users {
		userDTOS[i] = ToUserDTO(itm)
	}

	return userDTOS
}
