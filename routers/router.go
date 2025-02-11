package routers

import (
	"flight_api/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Define the API namespace
	ns := web.NewNamespace("/api/v1",
		web.NSNamespace("/flights",
			web.NSRouter("all_params/search", &controllers.FlightController{}, "get:GetByAllParams"),
      web.NSRouter("/:id", &controllers.FlightController{}, "get:GetFlightDetails"),
		),
	)

	// Register the namespace
	web.AddNamespace(ns)
  beego.Router("/swagger/*", &controllers.SwaggerController{})
}
