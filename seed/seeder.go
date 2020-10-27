package seed

import (
	"fmt"
	"github.com/simple-jwt-auth/middleware/auth"
	"github.com/simple-jwt-auth/models"
	"gorm.io/gorm"
	"log"
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
	}, {
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
	}, {
		Model:         gorm.Model{},
		UserName:      "admin",
		Password:      "bigdata@2019",
		Sub:           "f",
		Name:          "f",
		GivenName:     "f",
		FamilyName:    "f",
		Profile:       "f",
		Picture:       "f",
		Email:         "admin@dpbdhub.com",
		EmailVerified: true,
		Gender:        "1",
	},
}

func Load(db *gorm.DB) {
	dbExist := db.Migrator().HasTable(&models.User{})
	if !dbExist {
		err := db.Debug().Create(&models.User{}).Error
		if err != nil {
			log.Fatal("Migration error: %s", err.Error())
		}
	} else {
		db.Migrator().DropTable(&models.User{})
		db.Debug().Migrator().CreateTable(&models.User{})
	}

	//delete old data
	db.Debug().Model(&models.User{}).Delete(&models.User{}).Where("1 = 1")

	// sync new data
	for _, v := range users {
		err := db.Debug().Model(&models.User{}).Create(&v).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

	}

	//create default rbac rule
	enforcer := auth.NewCasbinEnforcerFromDB(db)
	//create default policy
	enforcer.AddPolicy("admin", "/auth/policy", "GET")
	enforcer.AddPolicy("admin", "/auth/policy", "POST")

	enforcer.AddPolicy("admin", "/jwt/auth/policy", "GET")
	enforcer.AddPolicy("admin", "/jwt/auth/policy", "POST")
	enforcer.AddPolicy("admin", "/jwt/auth/policy", "DELETE")
	////create default policy
	enforcer.AddPolicy("admin", "/jwt/auth/grouppolicy/*", "GET")
	enforcer.AddPolicy("admin", "/jwt/auth/grouppolicy", "POST")

	hasPermiss, _ := enforcer.Enforce("admin", "/jwt/auth/policy", "DELETE")
	st := fmt.Sprintf("%t", hasPermiss)
	fmt.Println(st)
	hasTable := db.Migrator().HasTable(&models.User{})
	if hasTable {
		db.Migrator().DropTable(&models.User{})
	}

	err := db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	for i, _ := range users {
		err := db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
