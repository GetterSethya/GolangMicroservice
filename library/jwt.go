package library

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(f AppHandler) AppHandler {

	return func(w http.ResponseWriter, r *http.Request) (int, error) {
		secret := os.Getenv("JWT_SECRET")
		log.Println(secret)

		// get Authorization header
		jwtToken := GetTokenFromRequest(r)

		//
		token, err := ValidateJWT(jwtToken, secret)
		if err != nil || !token.Valid {

			return http.StatusUnauthorized, fmt.Errorf("invalid token")
		}

		// call appHandler func
		if status, err := f(w, r); err != nil {

			return status, err
		}

		return http.StatusOK, nil
	}

}

// subject berisi userId
//
// issuer "gomedsos"
// expiration time default 12 hours, as time after the token issued + 12 hours
//
// not before: The "nbf" (not before) claim identifies the time before which the JWT
//
// issued At: The "iat" (issued at) claim identifies the time at which the JWT was issued.  This claim can be used to determine the age of the JWT.MUST NOT be accepted for processing
// time.Now().Add(time.Hour * 12)
func CreateJWT(userId, secret string, expiry time.Time) (string, error) {
	exp := jwt.NewNumericDate(expiry)
	nbf := jwt.NewNumericDate(time.Now())
	iat := jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		Issuer:    "gomedsos",
		ExpiresAt: exp,
		NotBefore: nbf,
		IssuedAt:  iat,
	})

	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error when signing accessToken %+v", err.Error())
		return "", err
	}

	return accessToken, nil
}

// hanya panggil fungsi ini di route yang sudah ada dijaga oleh middleware JWTMiddleware
func GetUserIdFromJWT(r *http.Request) string {
	secret := os.Getenv("JWT_SECRET")
	token := GetTokenFromRequest(r)
	jwtToken, err := ValidateJWT(token, secret)

	if err != nil {
		return ""
	}

	return jwtToken.Claims.(jwt.MapClaims)["sub"].(string)
}

func ValidateJWT(token, secret string) (*jwt.Token, error) {

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %+v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func GetTokenFromRequest(r *http.Request) string {
	jwtToken := r.Header.Get("Authorization")

	return jwtToken
}
