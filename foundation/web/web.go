// Package web contains a small web framework extension
package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
	"os"
	"syscall"
)

//Handler is a type that handles the http request within our own little mini framework.
type Handler func(context context.Context, w http.ResponseWriter, r *http.Request) error

//App is entrypoint into our application and configured our context object for each of our http
//handlers. Feel free to add any configuration data/logic on this App struct
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

//NewApp create an App value that handle the set of routes for the application
func NewApp(shutdown chan os.Signal) *App {
	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
	return &app
}

//SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

//Handle ...
func (a *App) Handle(method string, path string, readiness Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {

		//boilerplate

		if err := readiness(r.Context(), w, r); err != nil {
			a.SignalShutdown()
			return
		}

		//boilerplate

	}
	a.ContextMux.Handle(method, path, h)
}
