package models

type Enviroment struct {
	RedisConfig         RedisConf
	SqlConfig           SqlConf
	GoogleConf          Google
	FacebookConf        Facebook
	GithubConf          Github
	Port                string `json:"port"`
	CasbinWatcherEnable bool   `json:"casbin_watcher_enable"`
}

type RedisConf struct {
	Host     string `json:host`
	Port     string `json:port`
	Username string `json:username`
	Password string `json:password`
}

type SqlConf struct {
	Driver   string `json:"driver"`
	Username string `json:"username"`
	Passord  string `json:"password"`
	Host     string `json:host`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type Google struct {
	ClientID                string `json:"client_id"`
	ProjectID               string `json:"project_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUri string `json:"auth_provider_x_509_cert_uri"`
	ClientSecret            string `json:"client_secret"`
	RedirectUrl             string `json:"redirect_url"`
}

type Github struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUrl  string `json:"redirect_url"`
}

type Facebook struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUrl  string `json:"redirect_url"`
}
