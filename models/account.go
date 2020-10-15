package models

import (
	"time"
)

type Accounts []Account

type Account struct {
	ID            uint   `gorm:"primary_key"`
	AccName       string `gorm:"type:VARCHAR(255);"`
	AccOwner      string `gorm:"type:VARCHAR(255);"`
	AccURL        string `gorm:"type:VARBINARY(512);"`
	AccType       string `gorm:"type:VARBINARY(255);"`
	AccKey        string `gorm:"type:VARBINARY(255);"`
	AccUser       string `gorm:"type:VARBINARY(255);"`
	AccPass       string `gorm:"type:VARBINARY(255);"`
	AccError      string `gorm:"type:VARBINARY(512);"`

	CreatedAt     time.Time  `deepcopier:"skip"`
	UpdatedAt     time.Time  `deepcopier:"skip"`
	DeletedAt     *time.Time `deepcopier:"skip" sql:"index"`
}

