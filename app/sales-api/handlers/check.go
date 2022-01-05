package handlers

import (
	"context"
	"github.com/RymarSergey/ardanolabs2/foundation/web"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net/http"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return web.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}
	status := struct {
		Status string
	}{Status: "OK"}

	return web.Respond(ctx, w, status, http.StatusOK)
}
