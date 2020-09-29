package models

import (
	"errors"
	"log"
	"strconv"
)

var us = []User{
	{
		ID:       "2",
		UserName: "users",
		Password: "pass",
	}, {
		ID:       "3",
		UserName: "username",
		Password: "password",
	},
}
var UserRepo = UserRepository{
	Users: us,
}

type UserRepository struct {
	Users []User
}

func (r *UserRepository) FindAll() ([]User, error) {
	return r.Users, nil
}

func (r *UserRepository) FindByID(id int) (User, error) {

	for _, v := range r.Users {
		uid, err := strconv.Atoi(v.ID)
		if err != nil {
			return User{}, err
		}
		if uid == int(id) {
			return v, nil
		}
	}

	return User{}, errors.New("Not found")
}

func (r *UserRepository) Save(user User) (User, error) {
	r.Users = append(r.Users, user)

	return user, nil
}

func (r *UserRepository) Delete(user User) {
	id := -1
	for i, v := range r.Users {
		if v.ID == user.ID {
			id = i
			break
		}
	}

	if id == -1 {
		log.Fatal("Not found user ")
		return
	}

	r.Users[id] = r.Users[len(r.Users)-1] // Copy last element to index i.
	r.Users[len(r.Users)-1] = User{}      // Erase last element (write zero value).
	r.Users = r.Users[:len(r.Users)-1]    // Truncate slice.

	return
}
