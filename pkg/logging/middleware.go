package logging

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	n, err := w.ResponseWriter.Write(data)
	return n, err
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("New request to '%s' endpoint", r.Pattern)
		lrw := newLoggingResponseWriter(w)
		next(lrw, r)
		log.Printf("Request to endpoint '%s' processed. Status code: %d", r.Pattern, lrw.statusCode)

		if lrw.statusCode == http.StatusInternalServerError {
			log.Println(string(lrw.body))
			return
		}
	}
}
