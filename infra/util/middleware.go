package util

import (
	"log"
	"net/http"
	"strings"

	"github.com/changpro/disk-service/infra/config"
	"github.com/golang-jwt/jwt"
)

type Claim struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func AddMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Println("request origin is " + origin)
		if strings.Contains(origin, config.GetConfig().RequestOrigin) ||
			strings.Contains(origin, "localhost") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			h.ServeHTTP(w, r)
			return
		}
		u := r.URL.String()
		if strings.Contains(u, "sign-in") ||
			strings.Contains(u, "sign-up") ||
			strings.Contains(u, "modify-pw") ||
			strings.Contains(u, "file/download") {
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
