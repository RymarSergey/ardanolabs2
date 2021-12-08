// package handlers contains the full set of handlers and routes
// suported by the web api
package handlers

import (
	"github.com/dimfeld/httptreemux/v5"
	"log"
	"net/http"
	"os"
)

//API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger) *httptreemux.ContextMux {
	tm := httptreemux.NewContextMux()
	check := check{log: log}
	tm.Handle(http.MethodGet, "/test", check.readiness)

	return tm
}
