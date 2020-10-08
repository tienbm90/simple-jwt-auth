package models

type Enviroment struct {
	RedisConfig RedisConf
	SqlConfig   SqlConf
	
	Port string `json:"port"`
}

type RedisConf struct {
	Host     string `json:host`
	Port     string `json:port`
	Password string `json:password`
}

type SqlConf struct {
	Username string `json:"username"`
	Passord  string `json:"password"`
	Url      string `json:url`
}
