// package handlers contains the full set of handlers and routes
// suported by the web api
package handlers

import (
	"github.com/RymarSergey/ardanolabs2/business/mid"
	"github.com/RymarSergey/ardanolabs2/foundation/web"
	"log"
	"net/http"
	"os"
)

//API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log))

	check := check{
		log: log,
	}
	app.Handle(http.MethodGet, "/readiness", check.readiness)

	return app
}
