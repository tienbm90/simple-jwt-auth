package auth

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"log"
)

func NewCasbinEnforcer(connStr string) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapter("mysql", connStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Sqlerror: %s", err.Error()))
	}
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		log.Fatal(fmt.Sprintf("Creating enforcer error: %s", err.Error()))
	}

	////create default policy
	//enforcer.AddPolicy("admin","/auth/policy", "GET")
	//enforcer.AddPolicy("admin","/auth/policy", "POST")
	////create default policy
	//enforcer.AddPolicy("admin","/auth/grouppolicy/*", "GET")
	//enforcer.AddPolicy("admin","/auth/grouppolicy", "POST")
	//
	//enforcer.LoadPolicy()
	return enforcer
}
