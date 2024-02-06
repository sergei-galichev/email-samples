package app

import (
	"log"
)

type App struct {
	provider *serviceProvider
}

func NewApp() *App {
	app := &App{}
	app.initDeps()

	return app
}

func (a *App) Run() {
	a.provider.HttpServer().Run()
}

func (a *App) Stop() {
	log.Println("Stopping the App...")
	a.provider.storage.Close()
	close(a.provider.mailChan)
	log.Println("All resources are closed. The App is stopped")
}

func (a *App) initDeps() {
	inits := []func(){
		a.initServiceProvider,
		a.initHttpServer,
		a.initEmailServer,
	}

	for _, fn := range inits {
		fn()
	}
}

func (a *App) initServiceProvider() {
	a.provider = NewServiceProvider()
}

func (a *App) initHttpServer() {
	a.provider.HttpServer().Routes()
}

func (a *App) initEmailServer() {
	go a.provider.EmailServer().MailDelivery()
}
