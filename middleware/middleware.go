package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Replace this with your actual JWT secret key
var jwtKey = []byte("my_secret_key")

// Claims struct for JWT token claims.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWTMiddleware checks the JWT token before allowing access to the endpoint.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the request header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// Split the "Bearer <token>" format
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the token is valid, proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
