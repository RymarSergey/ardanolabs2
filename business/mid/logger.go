package mid

import (
	"context"
	"github.com/RymarSergey/ardanolabs2/foundation/web"
	"log"
	"net/http"
	"time"
)

//Logger write some information about to request to the log in the
//format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) web.Middleware {

	//This is the actual middleware function to be executed.
	m := func(beforeAfter web.Handler) web.Handler {

		//Create the handler that will be attached to the middleware chain.
		h := func(context context.Context, w http.ResponseWriter, r *http.Request) error {
			//If the context is missing the value, request the service
			//to be shutdown gracefully.
			v, ok := context.Value(web.KeyValue).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			log.Printf("%s : started   : %s %s -> %s",
				v.TraceID,
				r.Method, r.URL.Path, r.RemoteAddr,
			)

			err := beforeAfter(context, w, r)

			log.Printf("%s : completed : %s %s -> %s (%d) (%s)",
				v.TraceID,
				r.Method, r.URL.Path, r.RemoteAddr,
				v.StatusCode, time.Since(v.Now),
			)

			return err
		}
		return h
	}
	return m
}
