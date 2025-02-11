package controllers

import (
	"github.com/beego/beego/v2/server/web"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SwaggerController will serve the Swagger UI using the standard HTTP handler
type SwaggerController struct {
	web.Controller
}

func (c *SwaggerController) Get() {
	httpSwagger.Handler().ServeHTTP(c.Ctx.ResponseWriter, c.Ctx.Request)
}