package handlers

import (
	"boilerplate/api/authentication/generation"
	"boilerplate/config"
	"boilerplate/usecases/students"
	"boilerplate/usecases/universities"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	logger              *log.Logger
	config              config.Settings
	studentService      students.Service
	universitiesService universities.Service
}

func New(logger *log.Logger, config config.Settings, ss students.Service, us universities.Service) *Handler {
	return &Handler{
		logger:              logger,
		config:              config,
		studentService:      ss,
		universitiesService: us,
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
