package models

import "gorm.io/gorm"

type Policy struct {
	User   string `json:"user" forms:"user" query:"user"`
	Path   string `json:"path" forms:"path" query:"path"`
	Method string `json:"method" forms:"method" query:"method"`
}

type GroupPolicy struct {
	Member string `json:"member" forms:"member" query:"member"`
	Group  string `json:"group" forms:"group" query:"Group"`
}

type CasbinRule struct {
	gorm.Model
	PType string `gorm:"varchar(100) index not null default ''" json:"pType"`
	V0    string `gorm:"varchar(100) index not null default ''" json:"v0"`
	V1    string `gorm:"varchar(100) index not null default ''" json:"v1"`
	V2    string `gorm:"varchar(100) index not null default ''" json:"v2"`
	V3    string `gorm:"varchar(100) index not null default ''" json:"v3"`
	V4    string `gorm:"varchar(100) index not null default ''" json:"v4"`
	V5    string `gorm:"varchar(100) index not null default ''" json:"v5"`
}

func (rule *CasbinRule) TableName() string {
	return "casbin_rule"
}
