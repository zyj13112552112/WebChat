package main

import (
	"chat/config"
	"chat/mysql"
	"chat/router"
)

//项目入口

func main() {
	mysql.Migration()
	r := router.Router()
	r.Run(config.Port)
}
