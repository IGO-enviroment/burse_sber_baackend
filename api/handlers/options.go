package handlers

import "net/http"

func (s *Handler) Options(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
