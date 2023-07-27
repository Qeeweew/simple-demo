package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"simple-demo/common/config"
	"simple-demo/common/db"
	"simple-demo/common/log"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		config.ReadCfg()
		config.Init()
		log.Init()
		db.Init()
	})
}

func main() {
	r := gin.Default()

	initRouter(r)

	r.Run(fmt.Sprintf("%s:%s", config.AppCfg.Host, config.AppCfg.Port))
}
