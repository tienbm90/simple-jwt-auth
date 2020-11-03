package models

import (
	"gorm.io/gorm"
)

type OauthClientRepository struct {
	DB *gorm.DB
}

func ProvideOauthClientRepository(DB *gorm.DB) OauthClientRepository {
	return OauthClientRepository{DB: DB}
}

func (r *OauthClientRepository) FindAll() ([]OauthClient, error) {
	var client []OauthClient
	err := r.DB.Debug().Model(&OauthClient{}).Scan(&client).Error
	return client, err
}

func (r *OauthClientRepository) FindByID(id int) (OauthClient, error) {
	var client OauthClient
	err := r.DB.Debug().Model(&OauthClient{}).First(&client, id).Error
	return client, err
}

func (r *OauthClientRepository) Create(client OauthClient) (OauthClient, error) {
	res := r.DB.Debug().Model(&OauthClient{}).Create(&client)
	return client, res.Error
}

func (r *OauthClientRepository) Update(client OauthClient) (OauthClient, error) {
	res := r.DB.Debug().Model(&OauthClient{}).Updates(client)
	return client, res.Error
}

func (r *OauthClientRepository) Delete(client OauthClient) (OauthClient, error) {
	res := r.DB.Delete(&client)
	return client, res.Error
}
