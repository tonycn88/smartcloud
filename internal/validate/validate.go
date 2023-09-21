package validate

import (
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"smartcloud/internal/config"
)

type Validator interface {
	validate(r *http.Request) bool
	enable() bool
}

type ValidatorService struct {
	validators []Validator
}

func init_roots(caPath string) {
	// 创建自定义的证书池
	Roots = x509.NewCertPool()

	// 加载根证书
	caCertFile, err := os.ReadFile(caPath)
	if err != nil {
		log.Printf("无法读取根证书：%v", err)
		return
	}

	// 将根证书添加到证书池
	ok := Roots.AppendCertsFromPEM(caCertFile)
	if !ok {
		log.Println("无法添加根证书到证书池")
		return
	}

	opts = x509.VerifyOptions{
		Roots: Roots,
	}

}

func NewValidatorService(conf config.Authentication) *ValidatorService {

	var vs []Validator
	if conf.Type == "basic" {
		v := NewBasciValidator(conf.Basic.Enable)
		vs = append(vs, v)
	} else if conf.Type == "ca" {
		init_roots(conf.Ca.CaCrt)
		v := NewBx509Validator(conf.Ca.Enable)
		vs = append(vs, v)
	} else if conf.Type == "any" {
		v := NewBasciValidator(conf.Basic.Enable)
		vs = append(vs, v)
		x := NewBx509Validator(conf.Ca.Enable)
		vs = append(vs, x)
	}
	return &ValidatorService{
		validators: vs,
	}
}

func (v *ValidatorService) Validate(r *http.Request) bool {
	for _, validate := range v.validators {

		if validate.enable() && validate.validate(r) {
			return true
		}
	}
	return false
}
