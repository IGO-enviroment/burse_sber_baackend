package handlers

import (
	"boilerplate/gen"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginModel gen.Login
	if err := json.NewDecoder(r.Body).Decode(&loginModel); err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenResponse, err := s.studentService.Authenticate(loginModel)
	if err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(tokenResponse)
	if err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add(ContentTypeHeader, JsonContentType)
	w.Write(b)

	return
}
