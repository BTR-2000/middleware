package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return handler
	}

	total := len(middlewares) - 1
	for i := total; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func MustMiddlewares(handler http.Handler) http.Handler {
	slice := []Middleware{LoggingMiddleware, AuthMiddleware, RecoverMiddleware}
	return Chain(handler, slice...)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		time.Sleep(time.Millisecond)
		next.ServeHTTP(rw, r)

		log.Println("Метод:", r.Method)
		log.Println("Путь:", r.URL.Path)
		log.Println("Статус:", rw.statusCode)
		log.Printf("Времени потребовалось: %d ms.\n\n", time.Since(start).Milliseconds())
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "ПАНИКА", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "qwerty" {
			http.Error(w, "Не авторизовались", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if !strings.HasPrefix(contentType, "application/json") {
			http.Error(w, "Content-Type должен быть в JSON", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
