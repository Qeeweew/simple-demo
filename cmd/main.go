package main

import (
	"fmt"
	"simple-demo/common/config"
	"simple-demo/common/log"
	"simple-demo/repository"

	"github.com/gin-gonic/gin"
)

func init() {
	config.ReadCfg()
	config.Init()
	log.Init()
	repository.Init()
}

func main() {
	r := gin.Default()

	initRouter(r)

	r.Run(fmt.Sprintf(":%s", config.AppCfg.Port))
}
