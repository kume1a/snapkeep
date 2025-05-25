package shared

import (
	"log"
	"snapkeep/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	return generateJWT(env.AccessTokenSecret)
}

func VerifyAccessToken(tokenString string) error {
	env, err := config.ParseEnv()
	if err != nil {
		return err
	}

	return verifyJWT(tokenString, env.AccessTokenSecret)
}

func verifyJWT(tokenString string, secretKey string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("Error parsing token: ", err)
		return err
	}

	if !token.Valid {
		log.Println("Token is invalid")
		return err
	}

	return nil
}

func generateJWT(secretKey string) (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	expDuration := time.Second * time.Duration(env.AccessTokenExpSeconds)
	expiresAt := time.Now().Add(expDuration)

	claims := jwt.MapClaims{
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
		"sub": "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
