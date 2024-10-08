package routes

import (
	loginpage "sample/loginPage"
	"sample/transfer"
	"sample/web"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetupRoutes(app *fiber.App) {

	// PROMETHEUS
	prometheus := fiberprometheus.New("fiber_app")
	prometheus.RegisterAt(app, "/metrics")

	apiEndpoint := app.Group("/api")
	v1Endpoint := apiEndpoint.Group("/v1")

	app.Use(prometheus.Middleware)

	v1Endpoint.Post("/register", loginpage.Registration) //Registration
	v1Endpoint.Post("/log", loginpage.LoginPages)        //Login
	// app.Use(jwttoken.TokenAuthMiddleware)
	v1Endpoint.Post("/Ftransaction", transfer.ProcessCreditTransfer)

	// MONITOR
	monitorEndpoint := v1Endpoint.Group("/monitor")
	monitorEndpoint.Get("/", monitor.New(monitor.Config{
		Title:   "Monitor Instapay UAT",
		Refresh: 1,
	}))

	// WEB
	webEnpoint := v1Endpoint.Group("/web")
	webEnpoint.Get("/show-login", web.ShowLogin)
	webEnpoint.Get("/show-register", web.ShowRegister)
}
