package http_v1

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"mail-server/internal/model"
	services "mail-server/internal/service"
	"net/http"
	"os"
)

type HttpServer struct {
	mux *chi.Mux
	srv *http.Server

	mailService services.MailService
	mailChan    chan model.Mail
}

func NewServer(
	mailService services.MailService,
	mailChan chan model.Mail,
) *HttpServer {
	mux := chi.NewRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("PORT")),
		Handler: mux,
	}

	return &HttpServer{
		mux:         mux,
		srv:         srv,
		mailService: mailService,
		mailChan:    mailChan,
	}
}

func (s *HttpServer) Run() {
	log.Printf("HTTP Server is started on address: %s", s.srv.Addr)
	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("Shutting down the Mail App HttpServer")
		}
	}()
}
