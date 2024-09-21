package handlers

import (
	"boilerplate/config"
	"boilerplate/usecases/students"
	"log"
)

type Handler struct {
	logger         *log.Logger
	config         config.Settings
	studentService students.Service
}

func New(logger *log.Logger, config config.Settings, ss students.Service) *Handler {
	return &Handler{
		logger:         logger,
		config:         config,
		studentService: ss,
	}
}
