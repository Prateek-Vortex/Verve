package appcontext

import (
	"Verve/internal/configs/logger"
	restclient "Verve/internal/configs/restClient"
	"Verve/internal/database"
	"Verve/internal/repository"
	"Verve/internal/service"
	"context"
	"log/slog"
)

type AppContext struct {
	Logger          *slog.Logger
	VerveService    service.VerveService
	VerveRepository repository.VerveRepository
	RestClient      *restclient.RestClient
}

var appContext *AppContext

func LoadAppContext(db database.Service) {
	appContext.Logger = logger.InitLogger("text")
	appContext.RestClient = restclient.NewRestClient()
	appContext.VerveRepository = repository.NewImplVerveRepository(db)
	appContext.VerveService = service.NewImplVerveService(appContext.VerveRepository, appContext.RestClient, appContext.Logger)
	initBackgroundTasks()
}

func initBackgroundTasks() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go appContext.VerveService.LogUniqueCountEveryMinute(ctx)
}

func GetAppContext() *AppContext {
	return appContext
}
