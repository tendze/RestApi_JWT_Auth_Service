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

func GenerateToken(login, password string, getter USERGetter) (string, error) {
	_, err := getter.UserExists(login, password)
	if err != nil {
		if err == storage.ErrUserNotFound {
			return "", fmt.Errorf("user does not exists or incorrect login or password")
		}
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("invalid expiration duration: %w", err)
	}
	claims := &jwt.StandardClaims{
		Subject:   login,
		ExpiresAt: time.Now().Add(jwtExpirationTime).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing key: %w", err)
	}
	return signedToken, nil
}