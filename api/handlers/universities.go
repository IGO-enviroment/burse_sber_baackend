package handlers

import (
	"boilerplate/gen"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Handler) AddStudents(w http.ResponseWriter, r *http.Request) {
	claims, err := s.getTokenClaims(r)
	if err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !claims.IsUniversity {
		s.logger.Println("request from not university")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var requestBody gen.AddStudent
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.universitiesService.AddStudents(requestBody)
	if err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(response)
	if err != nil {
		s.logger.Println(fmt.Sprintf("%v", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add(ContentTypeHeader, JsonContentType)
	w.Write(b)

	return
}
