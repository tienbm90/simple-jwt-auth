package form

type Login struct {
	UserName string `json: "UserName"`
	Password string `json: "password"`
}

func (f Login) HasUserName() bool {
	return f.UserName != "" && len(f.UserName) <= 255
}

func (f Login) HasPassword() bool {
	return f.Password != "" && len(f.Password) <= 255
}
