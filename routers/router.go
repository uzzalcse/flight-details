package routers

import (
	"flight-details/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/api",
			beego.NSRouter("/:id", &controllers.FlightController{}, "get:GetFlightDetails"),
		),

	)

	beego.AddNamespace(ns)
	beego.Router("/swagger/*", &controllers.SwaggerController{})
}