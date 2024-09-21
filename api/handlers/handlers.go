package handlers

import (
	"boilerplate/api/authentication/generation"
	"boilerplate/config"
	"boilerplate/usecases/students"
	"log"
	"net/http"
	"strings"
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

func (s *Handler) getTokenClaims(r *http.Request) (generation.AccessTokenClaims, error) {
	tokenStr := r.Header.Get("Authorization")
	parts := strings.Split(tokenStr, " ")
	tokenStr = parts[1]
	token, err := generation.NewJWTToken(tokenStr, &s.config)
	if err != nil {
		return generation.AccessTokenClaims{}, err
	}
	claims := token.Claims()
	if !claims.IsUniversity {
		return generation.AccessTokenClaims{}, err
	}

	return claims, nil
}
