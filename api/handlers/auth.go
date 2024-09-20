package handlers

import (
	"boilerplate/gen"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
)

func (s *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginModel gen.Login
	if err := json.NewDecoder(r.Body).Decode(&loginModel); err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO Validate login model.

	claims := jwt.StandardClaims{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(signedToken))
}
