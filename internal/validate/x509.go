package validate

import (
	"crypto/x509"
	"log"
	"net/http"
	"time"
)

var Roots *x509.CertPool
var opts x509.VerifyOptions

type x509Validator struct {
	enabled bool
}

func NewBx509Validator(enabled bool) *x509Validator {
	return &x509Validator{
		enabled: enabled,
	}
}

func (t *x509Validator) enable() bool {
	return t.enabled
}

func (t *x509Validator) validate(r *http.Request) bool {
	// 验证证书链的可信任性

	// 检查是否有客户端证书
	for _, rawCert := range r.TLS.PeerCertificates {
		// 获取客户端证书信息
		// clientCert := r.TLS.PeerCertificates[0]
		x := rawCert.Raw
		// 验证客户端证书

		cert, err := x509.ParseCertificate(x)

		if err != nil {
			log.Printf("Parse Certificate fail: %v", err)
			continue
		}
		err = cert.VerifyHostname(r.Host)
		if err != nil {
			log.Printf("Verify Hostname fail: %v", err)
			continue
		}
		// 示例：验证证书是否已过期
		if cert.NotAfter.Before(time.Now()) {
			log.Printf("cert is expired!")
			continue
		}

		// 示例：验证证书的主题是否与期望的一致
		expectedSubject := "CN=Client"
		if cert.Subject.String() != expectedSubject {
			log.Printf("Certificate Subject is not matched!")
			continue
		}

		_, err = cert.Verify(opts)
		if err != nil {
			log.Printf("Certificate chain verify failed:%v", err)
			// return fmt.Errorf("证书链验证失败：%v", err)
			continue
		}

		log.Printf("Certificate Subject: %s, Issuer: %s", cert.Subject, cert.Issuer)
		return true
	}
	log.Println("Certificate  verify failed")
	return false
}
