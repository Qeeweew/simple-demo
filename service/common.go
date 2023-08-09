package service

import "fmt"

func genChatKey(x uint, y uint) string {
	if x > y {
		x, y = y, x
	}
	return fmt.Sprintf("%d_%d", x, y)
}

type MessageLatest struct {
	Message string `redis:"message"`
	Sender  uint   `redis:"sender"`
}
