package web

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

//Respond convert Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	//Set the status code for the request logger middleware.
	//If the context is missing this value, request the service
	// to be shutdown gracefully.
	v, ok := ctx.Value(KeyValue).(*Values)
	if !ok {
		return NewShutdownError("web value missing from context")
	}
	v.StatusCode = statusCode

	//If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	//convert the response value to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(statusCode) //TODO http: superfluous response.WriteHeader call from github.com/RymarSergey/ardanolabs2/foundation/web.Respond (response.go:35)

	if _, err = w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

//RespondError  sends an error response back  to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {

	//If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(ctx, w, er, webErr.Status); err != nil {
			return err
		}
	}

	//If not, the handler send any arbitrary error value so use 500.
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, w, er, http.StatusInternalServerError); err != nil {
		return err
	}

	return nil
}
