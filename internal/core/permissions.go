package core

import (
	"log"
	"smartcloud/internal/config"

	"github.com/casbin/casbin/v2"
)

type Permission struct {
	Enforcer *casbin.Enforcer
}

func NewPermission(conf config.Permission) *Permission {
	e, _ := casbin.NewEnforcer(conf.Model, conf.Policy)
	return &Permission{Enforcer: e}
}

func (e *Permission) Enforce(user string, path string, method string) bool {
	res, _ := e.Enforcer.Enforce(user, path, method)
	log.Println("result is ", res)
	return res
}
