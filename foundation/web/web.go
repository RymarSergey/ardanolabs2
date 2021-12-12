// Package web contains a small web framework extension
package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
	"net/http"
	"os"
	"syscall"
	"time"
)

//ctxKey represents the type of value for the ctxKey
type ctxKey int

//KeyValue is how request values are stored/retrieved
const KeyValue ctxKey = 1

//
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

//Handler is a type that handles the http request within our own little mini framework.
type Handler func(context context.Context, w http.ResponseWriter, r *http.Request) error

//App is entrypoint into our application and configured our context object for each of our http
//handlers. Feel free to add any configuration data/logic on this App struct
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

//NewApp create an App value that handle the set of routes for the application
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
	return &app
}

//SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

//Handle ...
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {

	//First wrap handler specific middleware around this handler
	handler = wrapMiddleware(mw, handler)

	//Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		//Set the context with required values to process the request
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValue, &v)

		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}

		//boilerplate

	}
	a.ContextMux.Handle(method, path, h)
}
