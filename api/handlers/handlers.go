package handlers

import (
	"boilerplate/config"
	"log"
)

type Handler struct {
	logger *log.Logger
	config config.Settings
}

func New(logger *log.Logger, config config.Settings) *Handler {
	return &Handler{
		logger: logger,
		config: config,
	}
}
