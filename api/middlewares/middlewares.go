package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

	"devtipmebackend/api/responses"
)

func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		next.ServeHTTP(w, r)
	})
}

func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = map[string]interface{}{"status": "failed", "message": "Missing authorization token"}

		var header = r.Header.Get("Authorization")
		header = strings.TrimSpace(header)

		if header == "" {
			responses.JSON(w, http.StatusForbidden, resp)
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			resp["status"] = "failed"
			resp["message"] = "Invalid token, please login"
			responses.JSON(w, http.StatusForbidden, resp)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), "userID", claims["userID"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
