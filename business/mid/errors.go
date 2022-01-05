package mid

import (
	"context"
	"github.com/RymarSergey/ardanolabs2/foundation/web"
	"log"
	"net/http"
)

//Errors handles errors coming out of the call chain. It detects normal
//application errors which are used to respond to the client in a uniform way.
//Unexpected errors (status>=500) are logged.
func Errors(log *log.Logger) web.Middleware {

	//This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {
		//Create the handler that will be attached to the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			//If the context is missing the value, request the service
			//to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValue).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			//Run the handler chain and catch any propagate error
			if err := handler(ctx, w, r); err != nil {

				//Log the error
				log.Printf("%s : ERROR     : %v", v.TraceID, err)

				//Respond to the error
				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				//If we receive the shutdown error we need to return it
				//back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			return nil
		}
		return h
	}
	return m
}
