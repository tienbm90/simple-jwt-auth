package models

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	DB *gorm.DB
}

func ProvideUserRepository(DB *gorm.DB) UserRepository {
	return UserRepository{DB: DB}
}

func (r *UserRepository) FindAll() ([]User, error) {
	var users []User
	err := r.DB.Debug().Model(&User{}).Scan(&users).Error

	return users, err
}

func (r *UserRepository) FindByID(id int) (User, error) {
	var user User
	err := r.DB.Debug().Model(&User{}).First(&user, id).Error
	return user, err
}

func (r *UserRepository) Validate(user User) (bool, error) {
	var us []User
	err := r.DB.Debug().Model(&User{}).Where("user_name = ?", user.UserName).Scan(&us).Error

	if err != nil {
		log.Printf("%s", err.Error())
		return false, err
	}

	if len(us) > 0 {
		u := us[0]
		if u.UserName == user.UserName && u.Password == user.Password {
			return true, nil
		}
	}

	return false, errors.New("Can't find any record matching with input")
}

func (r *UserRepository) FindByEmail(email string) (User, error) {
	var user User
	err := r.DB.Debug().Model(&User{}).Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *UserRepository) Create(user User) (User, error) {
	res := r.DB.Debug().Model(&User{}).Create(&user)
	return user, res.Error
}

func (r *UserRepository) Update(user User) (User, error) {
	res := r.DB.Debug().Model(&User{}).Updates(user)
	return user, res.Error
}

func (r *UserRepository) Delete(user User) (User, error) {
	res := r.DB.Delete(&user)
	return user, res.Error
}
