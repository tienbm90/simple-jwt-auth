package api

import (
	"github.com/go-playground/assert/v2"
	"github.com/simple-jwt-auth/models"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestToUser(t *testing.T) {
	var userDto = UserDTO{
		Username: "test1",
		Email:    "test1@gmail.com",
		Password: "123456a@",
	}
	user := ToUser(userDto)
	assert.Equal(t, userDto.Username, user.UserName)
}

func TestToUserDTO(t *testing.T) {
	var user = models.User{
		Model:    gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: nil,
		},
		UserName: "test1",
		Email:    "test1@gmail.com",
		Password: "qweryuiio",
	}

	userDto := ToUserDTO(user)

	assert.Equal(t, userDto.Password, user.Password)
}

func TestToUserDTOs(t *testing.T) {
	var users = []models.User{
		{
			Model:    gorm.Model{ID: 1},
			UserName: "test1",
			Email:    "test1@gmail.com",
			Password: "123456tq",
		}, {
			Model:    gorm.Model{ID: 2},
			UserName: "test2",
			Email:    "test2@outlook.com",
			Password: "fadfsafsdf",
		},
	}

	userDtos := ToUserDTOs(users)

	assert.Equal(t, len(users), len(userDtos))
	assert.Equal(t, userDtos[0].Username, users[0].UserName)
	assert.Equal(t, userDtos[0].Email, users[0].Email)
	assert.NotEqual(t, userDtos[0].Password, users[1].Password)

}
