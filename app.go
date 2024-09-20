package main

import (
	"boilerplate/api"
	"boilerplate/api/handlers"
	"boilerplate/config"
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type Application struct {
	MainCtx           context.Context
	MainCtxCancelFunc context.CancelFunc
	Server            *http.Server
	Logger            *log.Logger
}

func Create(settings config.Settings, logger *log.Logger) *Application {
	ctx, cancelFunc := context.WithCancel(context.Background())
	handler := handlers.New(
		logger,
		settings,
	)
	server := api.NewServer(ctx, settings, handler)

	return &Application{
		MainCtx:           ctx,
		MainCtxCancelFunc: cancelFunc,
		Server:            server,
		Logger:            logger,
	}
}

func (a *Application) Run() {
	go func() {
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Fatalf("listen http server err:  %v", err)
		}
	}()
}

func (a *Application) Shutdown() {
	a.Logger.Println("Shutting down application...")

	servShutDownCtx, serverShutDownCancelFunc := context.WithTimeout(a.MainCtx, 3*time.Second)
	defer serverShutDownCancelFunc()
	if err := a.Server.Shutdown(servShutDownCtx); err != nil {
		a.Logger.Fatalf("server shutdown err: %v", err)
	}

	a.MainCtxCancelFunc()
	a.Logger.Println("Application shutdown")
}
