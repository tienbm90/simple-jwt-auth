package models

import (
	"gorm.io/gorm"
)

type PolicyRepository struct {
	DB *gorm.DB
}

func ProvidePolicyRepository(DB *gorm.DB) UserRepository {
	return UserRepository{DB: DB}
}

func (r *PolicyRepository) FindAll() ([]CasbinRule, error) {
	var users []CasbinRule
	err := r.DB.Debug().Model(&CasbinRule{}).Scan(&users).Error
	return users, err
}

func (r *PolicyRepository) FindByID(id int) (CasbinRule, error) {
	var rule CasbinRule
	err := r.DB.Debug().Model(&CasbinRule{}).First(&rule, id).Error
	return rule, err
}

func (r *PolicyRepository) Create(rule CasbinRule) (CasbinRule, error) {
	res := r.DB.Debug().Model(&CasbinRule{}).Create(&rule)
	return rule, res.Error
}

func (r *PolicyRepository) Update(rule CasbinRule) (CasbinRule, error) {
	res := r.DB.Debug().Model(&CasbinRule{}).Updates(rule)
	return rule, res.Error
}

func (r *PolicyRepository) Delete(rule CasbinRule) (CasbinRule, error) {
	res := r.DB.Delete(&rule)
	return rule, res.Error
}
