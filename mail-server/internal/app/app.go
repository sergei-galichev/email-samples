package app

type App struct {
	provider *serviceProvider
}

func NewApp() *App {
	app := &App{}
	app.initDeps()

	return app
}

func (a *App) Run() {
	a.provider.Server().Run()
}

func (a *App) Stop() {
	a.provider.storage.Close()

}

func (a *App) initDeps() {
	inits := []func(){
		a.initServiceProvider,
		a.initServer,
	}

	for _, fn := range inits {
		fn()
	}
}

func (a *App) initServiceProvider() {
	a.provider = NewServiceProvider()
}

func (a *App) initServer() {
	a.provider.Server().Routes()
}
