package models

type Enviroment struct {
	RedisConfig RedisConf
	SqlConfig   SqlConf
	
	Port string `json:"port"`
	CasbinWatcherEnable bool `json:"casbin_watcher_enable"`
}

type RedisConf struct {
	Host     string `json:host`
	Port     string `json:port`
	Username string `json:username`
	Password string `json:password`
}

type SqlConf struct {
	Username string `json:"username"`
	Passord  string `json:"password"`
	Url      string `json:url`
}
