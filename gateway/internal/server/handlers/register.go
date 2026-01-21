package register_handler

import (
	"net/http"
)

func RegisterHandler() http.HandlerFunc {
	type RegisterRequest struct {
		username string
		password string
	}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}
