package utils

import (
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var SecretKey = []byte("1111222233334444")

func CreateToken(userID uint) string {
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": fmt.Sprintf("%v", userID),
	}).SignedString(SecretKey)
	if err != nil {
		logrus.Panic("JWT creating error", err)
	}
	return t
}

func ParseToken(tokenString string) (res uint, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return
	}
	var aud jwt.ClaimStrings
	aud, _ = token.Claims.GetAudience()
	id, _ := strconv.Atoi(aud[0])
	res = uint(id)
	return
}
