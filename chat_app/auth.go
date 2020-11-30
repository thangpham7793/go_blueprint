package main

import "net/http"

type withAuthHandler struct {
	next http.Handler
}

func (h *withAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r)
}

//wrapper for any handler
func mustAuth(handler http.Handler) http.Handler {
	return &withAuthHandler{handler}
}
