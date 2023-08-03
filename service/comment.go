package service

import (
	"simple-demo/common/model"
	"sync"
)

type commentService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	commentInstance *commentService
	commentOnce     sync.Once
)
