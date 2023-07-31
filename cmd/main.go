package main

import (
	"fmt"
	"simple-demo/common/config"
	"simple-demo/common/log"
	"simple-demo/repository/dbcore"

	"github.com/gin-gonic/gin"
)

func init() {
	config.ReadCfg()
	config.Init()
	log.Init()
	dbcore.Init()
}

func main() {
	r := gin.Default()

	initRouter(r)

	r.Run(fmt.Sprintf(":%s", config.AppCfg.Port))
}
