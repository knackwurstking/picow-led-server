package main

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	http.Hijacker

	StatusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		Hijacker:       w.(http.Hijacker),
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
