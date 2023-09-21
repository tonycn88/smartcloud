package validate

import "net/http"

type TokenValidator struct {
	enabled bool
}

func (t *TokenValidator) validate(r *http.Request) bool {
	return false
}
