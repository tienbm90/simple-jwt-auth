package seed

import (
	"github.com/simple-jwt-auth/models"
	"gorm.io/gorm"
)

var users = []models.User{
	{
		Model:         gorm.Model{},
		UserName:      "tien",
		Password:      "bigdata@2019",
		Sub:           "f",
		Name:          "f",
		GivenName:     "f",
		FamilyName:    "f",
		Profile:       "f",
		Picture:       "f",
		Email:         "tienbm90@gmail.com",
		EmailVerified: false,
		Gender:        "1",
	},{
		Model:         gorm.Model{},
		UserName:      "blackpresident",
		Password:      "bigdata@2019",
		Sub:           "f",
		Name:          "f",
		GivenName:     "f",
		FamilyName:    "f",
		Profile:       "f",
		Picture:       "f",
		Email:         "blackpresident90@gmail.com",
		EmailVerified: false,
		Gender:        "1",
	},
}

func Load(db *gorm.DB) {
	//err := db.Debug().AutoMigrate(&models.User{}).Error
	//if err != nil {
	//	log.Fatalf("cannot migrate table: %v", err)
	//}
	//for i, _ := range users {
	//	err := db.Debug().Model(&models.User{}).Create(&users[i]).Error
	//	if err != nil {
	//		log.Fatalf("cannot seed users table: %v", err)
	//	}
	//}
}
