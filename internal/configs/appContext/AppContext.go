package appcontext

import (
	"Verve/internal/configs/logger"
	restclient "Verve/internal/configs/restClient"
	"Verve/internal/database"
	"Verve/internal/event"
	"Verve/internal/repository"
	"Verve/internal/service"
	"context"
	"fmt"
	"log/slog"
)

type AppContext struct {
	Logger          *slog.Logger
	VerveService    service.VerveService
	VerveRepository repository.VerveRepository
	RestClient      *restclient.RestClient
	Event           event.Event
}

var appContext *AppContext

func LoadAppContext(db database.Service) {
	if appContext == nil {
		appContext = &AppContext{}
	}
	appContext.Logger = logger.InitLogger("text")
	appContext.RestClient = restclient.NewRestClient()
	event, err := event.NewKafkaEvent(appContext.Logger)
	appContext.Event = event
	if err != nil {
		appContext.Logger.Error("Failed to create kafka event", "error", err.Error())
		errs := fmt.Errorf("failed to create kafka event in load app context %w", err)
		panic(errs)
	}
	appContext.VerveRepository = repository.NewImplVerveRepository(db)
	appContext.VerveService = service.NewImplVerveService(appContext.VerveRepository, appContext.RestClient, appContext.Logger, appContext.Event)

	initBackgroundTasks()
}

func initBackgroundTasks() {
	go appContext.VerveService.LogUniqueCountEveryMinute(context.Background())
	go appContext.VerveService.SendUniqueCountEveryMinute(context.Background())
}

func GetAppContext() *AppContext {
	return appContext
}
