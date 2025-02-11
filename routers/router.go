package routers

import (
	"flight-details/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Define the API namespace
	ns := web.NewNamespace("/v1/api",
		web.NSNamespace("/flights",
			web.NSRouter("all_params/search", &controllers.FlightController{}, "get:GetByAllParams"),
			web.NSRouter("/:id", &controllers.FlightController{}, "get:GetFlightDetails"),
      web.NSRouter("/dest_time/search", &controllers.FlightController{}, "get:Get"), 
		),
	)

	// Register the namespace
	web.AddNamespace(ns)
	web.Router("/swagger/*", &controllers.SwaggerController{})
}
