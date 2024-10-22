package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"jwt-auth-service/internal/storage"
	"os"
	"time"
)

var secretKey = getSecretKey()
var jwtExpirationTime = getJwtExpirationTime()

func getSecretKey() []byte {
	_ = godotenv.Load()
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}
func getJwtExpirationTime() time.Duration {
	_ = godotenv.Load()
	expiration, _ := time.ParseDuration(os.Getenv("JWT_EXPIRATION_TIME"))
	return expiration
}

type USERGetter interface {
	UserExists(login, password string) (bool, error)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserLogin string `json:"login"`
}

func GenerateToken(login, password string, getter USERGetter) (string, error) {
	_, err := getter.UserExists(login, password)
	if err != nil {
		if err == storage.ErrUserNotFound {
			return "", fmt.Errorf("user does not exists or incorrect login or password")
		}
		return "", err
	}
	claims := &tokenClaims{jwt.StandardClaims{
		ExpiresAt: time.Now().Add(jwtExpirationTime).Unix(),
		IssuedAt:  time.Now().Unix(),
	}, login}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing key: %w", err)
	}
	return signedToken, nil
}

// ValidateToken returns user login and error
func ValidateToken(jwtToken string) (string, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if tok, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return tok.UserLogin, nil
	}
	return "", fmt.Errorf("invalid token")
}
