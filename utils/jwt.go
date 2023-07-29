package utils

import (
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// TODO: jwt鉴权

var secretKey = []byte("1111222233334444")

func CreateToken(userID uint) string {
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": fmt.Sprintf("%v", userID),
	}).SignedString(secretKey)
	if err != nil {
		logrus.Panic("JWT creating error", err)
	}
	return t
}

func ParseToken(tokenString string) (ans uint, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return
	}
	var aud jwt.ClaimStrings
	aud, err = token.Claims.GetAudience()
	if err != nil {
		return
	}
	id, err := strconv.Atoi(aud[0])
	if err != nil {
		return
	}
	ans = uint(id)
	return
}
