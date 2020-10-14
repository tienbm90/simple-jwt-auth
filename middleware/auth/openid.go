package auth

import (
	"crypto/rand"
	"encoding/base64"
)

type Credentials struct {
	Cid         string `json:"client_id"`
	Pid         string `json:"project_id"`
	Auri        string `json:auth_uri`
	Turi        string `json:"token_uri"`
	ProviderUri string `json:"auth_provider_x509_cert_url"`
	Csec        string `json:client_secret`
}


// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

//func getLoginURL(state string) string {
//	return Conf.AuthCodeURL(state)
//}
