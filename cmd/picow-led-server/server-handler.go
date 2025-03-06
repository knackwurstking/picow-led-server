package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type serverHandler struct{}

func (*serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered", "r", r)
		}
	}()

	rw := NewResponseWriter(w)

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	http.DefaultServeMux.ServeHTTP(rw, r)

	log := slog.Warn

	if rw.StatusCode >= 200 && rw.StatusCode < 300 || rw.StatusCode == 0 {
		log = slog.Info
	} else if rw.StatusCode >= 500 {
		log = slog.Error
	}

	log(fmt.Sprintf("%s %s", r.Method, r.URL.Path), "addr", r.RemoteAddr, "status", rw.StatusCode)
}
