package utils_test

import (
	"fmt"
	"simple-demo/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRightToken(t *testing.T) {
	token := utils.CreateToken(123)
	fmt.Println("token = ", token)
	id, err := utils.ParseToken(token)
	if err != nil {
		t.Errorf("wrong token")
		return
	}
	fmt.Println("id = ", id)
	if id != 123 {
		t.Errorf("worng id")
	}
}

func TestWrongToken(t *testing.T) {
	token := utils.CreateToken(123)
	token += "ad"
	fmt.Println("token = ", token)
	_, err := utils.ParseToken(token)
	assert.NotEqual(t, err, nil, "should not be nil")
}
