package validate

import (
	"log"
	"net/http"
	"smartcloud/internal/models"
)

type BasciValidator struct {
	enabled bool
}

func NewBasciValidator(enabled bool) *BasciValidator {
	return &BasciValidator{
		enabled: enabled,
	}
}

func (t *BasciValidator) enable() bool {
	return t.enabled
}
func (t *BasciValidator) validate(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	log.Printf("basic auth %s,%s", username, password)
	if !ok {
		return false
	}
	var u models.Users
	if !u.Exist(username) {
		log.Printf("%s not exists", username)
		return false
	}
	pass := u.Find_password_by_username(username)

	// 验证用户名/密码
	if password != pass {
		log.Printf("pass is %s,datapass is %s", password, pass)
		return false
	}

	return true
}
