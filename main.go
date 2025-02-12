package main

import (
	"flight-details/controllers"
	_ "flight-details/docs"
	_ "flight-details/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	controllers.Init()
	beego.Run()
}
