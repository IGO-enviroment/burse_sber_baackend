package main

import (
	"boilerplate/api"
	"boilerplate/api/handlers"
	"boilerplate/config"
	"boilerplate/postgres"
	"boilerplate/usecases/students"
	"boilerplate/usecases/universities"
	"context"
	"errors"
	mx "github.com/gorilla/handlers"
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
	pgDb, err := postgres.NewPostgresConnector(settings.PgConnString).Open()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	studentService := students.NewStudentsService(pgDb, settings)
	universitiesService := universities.NewUniversitiesService(pgDb, settings)
	handler := handlers.New(
		logger,
		settings,
		studentService,
		universitiesService,
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
		if err := http.ListenAndServe(":8080", mx.CORS()(a.Server.Handler)); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
