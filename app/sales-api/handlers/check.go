package handlers

import (
	"context"
	"github.com/RymarSergey/ardanolabs2/foundation/web"
	"log"
	"net/http"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	status := struct {
		Status string
	}{Status: "OK"}

	log.Println(r, status)
	return web.Respond(ctx, w, status, http.StatusOK)
}
