package main

import (
	"github.com/joho/godotenv"
	"log"
	"mail-server/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Config: Error loading .env file")
	}
}
func main() {
	a := app.NewApp()
	a.Run()

	// Shutdown
	shutdown([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM})
	a.Stop()
}

func shutdown(signals []os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	sig := <-ch
	log.Printf("Caught signal: %s. Shutting down...", sig)
}
