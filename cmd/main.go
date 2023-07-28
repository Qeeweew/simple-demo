package main

import (
	"fmt"
	"simple-demo/common/config"
	"simple-demo/common/db"
	"simple-demo/common/log"

	"github.com/gin-gonic/gin"
)

func init() {
	config.ReadCfg()
	config.Init()
	log.Init()
	db.Init()
}

func main() {
	r := gin.Default()

	initRouter(r)

	r.Run(fmt.Sprintf("%s:%s", config.AppCfg.Host, config.AppCfg.Port))
}
