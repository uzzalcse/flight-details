package routers

import (
    "flight_api/controllers"

    "github.com/beego/beego/v2/server/web"
)

func init() {
    // Initialize the Elasticsearch client (called once on app startup)
    controllers.Init()

    // Define the API namespace

    ns := web.NewNamespace("/api/v1",
        web.NSNamespace("/flights",
            web.NSRouter("/dest_time/search", &controllers.FlightController{}, "get:Get"),  // Corrected route
        ),
    )

    // Register the namespace
    web.AddNamespace(ns)
	web.Router("/swagger/*", &controllers.SwaggerController{})
}