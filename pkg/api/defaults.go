package api

import (
	"net/http"

	"github.com/DataDog/datadog-agent/pkg/api/util"
)

var LocalhostHosts = []string{"127.0.0.1", "localhost"}

func DefaultTokenValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := util.Validate(w, r); err != nil {
			return
		}
		next.ServeHTTP(w, r)
	})
}
