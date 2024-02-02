package http_v1

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (s *HttpServer) Routes() {
	s.mux.Use(middleware.Heartbeat("/api/ping"))
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)

	s.mux.Get("/", s.Home())
	s.mux.Post("/api/v1/submit", s.GetSubscriber())
	s.mux.Post("/api/v1/send", s.SendMail())

	fileServer := http.FileServer(http.Dir("./static"))
	s.mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

}
