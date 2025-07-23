package jwt

import (
	"time"

	"github.com/Alfian57/belajar-golang/internal/config"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/golang-jwt/jwt/v5"
	golangJwt "github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user model.User) (string, error) {

	token := golangJwt.NewWithClaims(golangJwt.SigningMethodHS256, golangJwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	secret := config.GetEnv("ACCESS_TOKEN_SECRET", "secret")
	secretByte := []byte(secret)

	tokenString, err := token.SignedString(secretByte)
	return tokenString, err
}

func CreateRefreshToken(user model.User) (string, error) {

	token := golangJwt.NewWithClaims(golangJwt.SigningMethodHS256, golangJwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	secret := config.GetEnv("REFRESH_TOKEN_SECRET", "secret")
	secretByte := []byte(secret)

	tokenString, err := token.SignedString(secretByte)
	return tokenString, err
}

func ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *golangJwt.Token) (any, error) {
		secret := config.GetEnv("ACCESS_TOKEN_SECRET", "secret")
		secretByte := []byte(secret)
		return secretByte, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(golangJwt.MapClaims); ok {
		return claims["id"].(string), nil
	}

	return "", err
}

func GetUserID(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *golangJwt.Token) (any, error) {
		secret := config.GetEnv("ACCESS_TOKEN_SECRET", "secret")
		secretByte := []byte(secret)
		return secretByte, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(golangJwt.MapClaims); ok {
		return claims["id"].(string), nil
	}

	return "", err
}
