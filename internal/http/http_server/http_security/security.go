package http_security

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var expTime = time.Now().Add(15 * (24 * time.Hour))

type JWTClaims struct {
	ID uuid.UUID `json:"ID"`
	jwt.StandardClaims
}

func DeleteJWTCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:    "jwt",
		Expires: time.Now().Add(-time.Hour),
	})
}

func GetJWTCookie(c echo.Context, secret []byte) (*JWTClaims, error) {
	var op = "/http/http_server/security.GetJWTCookie"

	cookie, err := c.Cookie("jwt")
	if err != nil {
		fmt.Println(op, ":", err)
		return nil, err
	}

	return validateJWT(cookie.Value, secret)
}

func SaveJWTCookie(c echo.Context, id uuid.UUID, secret []byte) error {
	var op = "/http/http_server/security.SaveJWT"
	token, err := generateJWT(id, secret)
	if err != nil {
		fmt.Println(op, ":", err)
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		Expires:  expTime,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func generateJWT(id uuid.UUID, secret []byte) (string, error) {
	var op = "/http/http_server/security.GenerateJWT"

	claims := &JWTClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			Issuer:    "localhost",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign, err := token.SignedString(secret)
	if err != nil {
		fmt.Println(op, ":", err)
		return "", err
	}

	return sign, nil
}

func validateJWT(str string, secret []byte) (*JWTClaims, error) {
	var op = "/http/http_server/security.ValidateJWT"
	token, err := jwt.ParseWithClaims(str, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		fmt.Println(op, ":", err)
		return nil, err
	}

	if token.Valid {
		claims, ok := token.Claims.(*JWTClaims)
		if ok {
			return claims, nil
		}
	}
	return nil, fmt.Errorf("invalid JWT token! ")
}
