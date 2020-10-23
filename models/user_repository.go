package models

import (
	"gorm.io/gorm"
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

func (r *UserRepository) Validate(user User) (User, error) {
	var us []User
	//res := r.DB.Debug().Model(&User{}).Where("email = ?", user.Email).Scan(&us)
	err := r.DB.Debug().Model(&User{}).Where("email = ?", user.Email).Scan(&us).Error

	u := us[0]
	if err != nil {
		return User{}, err
	}

	if u.UserName == user.UserName && u.Password == user.Password {
		return u, nil
	}

	return User{}, err
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
