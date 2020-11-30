package main

import (
	"chat_app/trace"
	"io"
	"net/http"
)

type logHandler struct {
	tracer trace.Tracer
	next   http.Handler
}

func (h *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.tracer.Trace(r.Method, r.URL)
	h.next.ServeHTTP(w, r)
}

func withLogger(handler http.Handler, w io.Writer) http.Handler {
	return &logHandler{next: handler, tracer: trace.New(w)}
}
