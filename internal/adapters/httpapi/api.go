package httpapi

import (
	"context"
	"github.com/xakpro/sms/internal/domain/core"
	"github.com/xakpro/sms/internal/interfaces"
	"net/http"
	"time"
)

type API struct {
	lg interfaces.Logger
	cr *core.Core

	server *http.Server
	lChan  chan error
}

func CreateAPI(lg interfaces.Logger, listen string, cr *core.Core) *API {
	api := &API{
		lg:    lg,
		cr:    cr,
		lChan: make(chan error, 1),
	}

	api.server = &http.Server{
		Addr:              listen,
		Handler:           api.router(),
		ReadTimeout:       2 * time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return api
}

func (a *API) Start() {
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.lg.Errorw("Http server closed", err)
			a.lChan <- err
		}
	}()
}

func (a *API) Wait() <-chan error {
	return a.lChan
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
