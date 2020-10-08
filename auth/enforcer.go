package auth

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"log"
)


func NewCasbinEnforcer(connStr string) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapter("mysql", connStr)
	if err != nil {
		log.Fatal("Sqlerror: %s", err.Error())
	}
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		log.Fatal("Creating enforcer error: %s", err.Error())
	}

	enforcer.AddPolicy()

	return enforcer
}
