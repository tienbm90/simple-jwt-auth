package auth

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"log"
)

func NewCasbinEnforcer(connStr string) *casbin.Enforcer {

	adapter, err := gormadapter.NewAdapter("mysql", connStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Sqlerror: %s", err.Error()))
	}

	rbac_model_with_resources_roles_model := `
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	g2 = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act
	`

	casbin_model, _ := model.NewModelFromString(rbac_model_with_resources_roles_model)

	enforcer, err := casbin.NewEnforcer(casbin_model, adapter)
	if err != nil {
		log.Fatal(fmt.Sprintf("Creating enforcer error: %s", err.Error()))
	}

	////create default policy
	enforcer.AddPolicy("admin", "/auth/policy", "GET")
	enforcer.AddPolicy("admin", "/auth/policy", "POST")
	////create default policy
	enforcer.AddPolicy("admin", "/auth/grouppolicy/*", "GET")
	enforcer.AddPolicy("admin", "/auth/grouppolicy", "POST")
	//
	//enforcer.LoadPolicy()
	return enforcer
}

func NewCasbinEnforcerFromDB(db *gorm.DB) *casbin.Enforcer {

	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "casbin", "rule")
	if err != nil {
		log.Fatal(fmt.Sprintf("Sqlerror: %s", err.Error()))
	}

	rbac_model_with_resources_roles_model := `
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	g2 = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act
	`

	casbin_model, _ := model.NewModelFromString(rbac_model_with_resources_roles_model)

	enforcer, err := casbin.NewEnforcer(casbin_model, adapter)
	if err != nil {
		log.Fatal(fmt.Sprintf("Creating enforcer error: %s", err.Error()))
	}

	////create default policy
	enforcer.AddPolicy("admin", "/auth/policy", "GET")
	enforcer.AddPolicy("admin", "/auth/policy", "POST")
	////create default policy
	enforcer.AddPolicy("admin", "/auth/grouppolicy/*", "GET")
	enforcer.AddPolicy("admin", "/auth/grouppolicy", "POST")
	//
	//enforcer.LoadPolicy()
	return enforcer
}
