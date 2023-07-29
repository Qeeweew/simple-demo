package utils

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// TODO: jwt鉴权

var secretKey = []byte("1111222233334444")

func createToken(userID uint) string {
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": userID,
	}).SignedString(secretKey)
	if err != nil {
		logrus.Panic("JWT creating error", err)
	}
	return t
}

func parseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return 0, errors.New("invalid token!")
	}
	aud, err := token.Claims.GetAudience()
	if err != nil {
		return 0, errors.New("invalid token!")
	}
	data, err := aud.MarshalJSON()
	if err != nil {
		return 0, errors.New("invalid token!")
	}
	id, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, errors.New("invalid token!")
	}
	return uint(id), nil
}
