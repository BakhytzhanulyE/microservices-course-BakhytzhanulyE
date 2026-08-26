package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/handler"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/middleware"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/storage"
)

// Адрес и таймауты держим константами, чтобы не разбрасывать «магические» числа по коду.
const (
	addr              = ":8080"
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	s := storage.NewStorage()
	h := handler.NewHandler(s)

	router := chi.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.Recoverer)

	router.Get("/health", h.Health)
	router.Get("/boom", func(_ http.ResponseWriter, _ *http.Request) {
		panic("boom!")
	})
	router.Get("/api/v1/orders/{order_uuid}", h.GetOrder)

	router.Post("/api/v1/orders", h.CreateOrder)
	router.Post("/api/v1/orders/{order_uuid}/pay", h.PayOrder)

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	log.Printf("сервер слушает на %s", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("сервер остановлен: %v", err)
	}
}
