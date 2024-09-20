package main

import (
	"boilerplate/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New(os.Stdout, "boilerplate", 0)
	settings, err := config.Read()
	if err != nil {
		logger.Fatalf("can't read settings: %v", err)
	}

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	app := Create(settings, logger)
	app.Run()

	<-stopCh
	app.Shutdown()

}
