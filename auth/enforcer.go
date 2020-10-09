package auth

import (
	"fmt"
	"github.com/billcobbler/casbin-redis-watcher"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	persist2 "github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/simple-jwt-auth/models"
	"log"
)

func NewCasbinEnforcer(connStr string) *casbin.Enforcer {

	adapter, err := gormadapter.NewAdapter("mysql", connStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Sqlerror: %s", err.Error()))
	}

	//rbac_model := `
	//
	//[request_definition]
	//r = sub, obj, act
	//
	//[policy_definition]
	//p = sub, obj, act
	//
	//[policy_effect]
	//e = some(where (p.eft == allow))
	//
	//[matchers]
	//m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
	//`

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

func NewCasbinEnforcerWithWatcher(env models.Enviroment) *casbin.Enforcer {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/", env.SqlConfig.Username, env.SqlConfig.Passord, env.SqlConfig.Url)

	adapter, err := gormadapter.NewAdapter("mysql", dataSource)
	if err != nil {
		log.Fatal(fmt.Sprintf("Sqlerror: %s", err.Error()))
	}

	//rbac_model := `
	//
	//[request_definition]
	//r = sub, obj, act
	//
	//[policy_definition]
	//p = sub, obj, act
	//
	//[policy_effect]
	//e = some(where (p.eft == allow))
	//
	//[matchers]
	//m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
	//`

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

	if env.CasbinWatcherEnable {
		var watcher persist2.Watcher
		if env.RedisConfig.Username != "" && env.RedisConfig.Password != "" {
			w, _ := rediswatcher.NewWatcher(fmt.Sprintf("redis://%s:%s@%s:%s", env.RedisConfig.Username, env.RedisConfig.Password, env.RedisConfig.Host, env.RedisConfig.Port))
			watcher = w
		} else {
			w, _ := rediswatcher.NewWatcher(fmt.Sprintf("%s:%s", env.RedisConfig.Host, env.RedisConfig.Port))
			watcher = w
		}

		enforcer.SetWatcher(watcher)
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
