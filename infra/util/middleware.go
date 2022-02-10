package util

import (
	"net/http"
	"strings"

	"github.com/changpro/disk-service/infra/config"
	"github.com/golang-jwt/jwt"
)

type Claim struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

// handle cors
func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Origin"), "http://localhost") {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "sign-in") {
			h.ServeHTTP(w, r)
			return
		}
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (i interface{}, err error) {
			return []byte(config.GetConfig().AuthKey), nil
		})
		if err != nil {
			var message string
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					message = "token is malformed"
				} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
					message = "token could not be verified because of signing problems"
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					message = "signature validation failed"
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					message = "token is expired"
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					message = "token is not yet valid before sometime"
				} else {
					message = "can not handle this token"
				}
			}
			http.Error(w, message, http.StatusUnauthorized)
			return
		}
		if _, ok := token.Claims.(*Claim); ok && token.Valid {
			h.ServeHTTP(w, r)
			return
		}
		http.Error(w, "cannot extract user from header", http.StatusUnauthorized)
	})
}
