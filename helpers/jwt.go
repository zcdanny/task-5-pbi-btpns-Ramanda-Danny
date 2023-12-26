package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/config"
)


func GenerateToken(userID uuid.UUID) (string, error) {
	secretKey := []byte(GetSecretKey())

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetSecretKey() string {
	secretKey := config.GetSecretKey()
	if secretKey == "" {
		secretKey = "defaultsecret"
	}
	return secretKey
}