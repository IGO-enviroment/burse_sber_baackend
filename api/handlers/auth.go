package handlers

import (
	"boilerplate/api/authentication/generation"
	"boilerplate/gen"
	"boilerplate/jwt"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginModel gen.Login
	if err := json.NewDecoder(r.Body).Decode(&loginModel); err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO Validate login model.

	accessTokenClaims := generation.AccessTokenClaims{
		UserId:              "user.UserId",
		UserName:            "user.UserName",
		Email:               "user.Email",
		AccountId:           0,
		IsEmployer:          false,
		IsBackofficeManager: false,
		IsStudent:           false,
		CreationTimestamp:   time.Now().UTC().Unix(),
		TTL:                 s.config.AccessTokenTTL,
	}

	token := jwt.GetToken(accessTokenClaims, s.config.SecretKey)

	w.Write([]byte(token))
}
