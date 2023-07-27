package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/oyaro-tech/auth/pkg/logger"
)

var JWT_SECRET []byte

func init() {
	godotenv.Load()

	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

	token, _ := generateToken("admin")
	logger.Debug.Println(token)
}

func generateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
