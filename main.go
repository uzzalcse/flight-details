package main

import (
	_ "flight-details/docs"
	_ "flight-details/routers"
	"flight-details/utils"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	utils.Init()
	beego.Run()
}
