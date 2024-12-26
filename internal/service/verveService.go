package service

import (
	restclient "Verve/internal/configs/restClient"
	"Verve/internal/model/entity"
	"Verve/internal/model/request"
	"Verve/internal/repository"
	"context"
	"log/slog"
	"time"
)

type VerveService interface {
	SaveAndPost(ctx context.Context, verveRequest request.VerveRequest) error
	LogUniqueCountEveryMinute(ctx context.Context)
}

type implVerveService struct {
	verveRepo  repository.VerveRepository
	restClient *restclient.RestClient
	Logger     *slog.Logger
}

func NewImplVerveService(repository repository.VerveRepository, client *restclient.RestClient, logger *slog.Logger) *implVerveService {
	return &implVerveService{
		verveRepo:  repository,
		restClient: client,
		Logger:     logger,
	}
}

func (vs *implVerveService) SaveAndPost(ctx context.Context, verveRequest request.VerveRequest) error {
	entity := entity.GetEntityFromRequest(verveRequest)
	err := vs.verveRepo.Save(ctx, entity)
	if err != nil {
		return err
	}
	if verveRequest.Url != "" {
		go vs.postToUrl(context.Background(), verveRequest.Url)
	}
	return nil
}

func (vs *implVerveService) postToUrl(ctx context.Context, url string) {
	count, err := vs.verveRepo.GetUniqueCount(ctx)
	if err != nil {
		vs.Logger.Error("Failed to get count", err)
	}
	jsonData := map[string]int64{"count": count}
	vs.restClient.Post(url, jsonData)
}

func (vs *implVerveService) LogUniqueCountEveryMinute(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				count, err := vs.verveRepo.GetUniqueCount(ctx)
				if err != nil {
					vs.Logger.Error("Failed to get unique count", "error", err)
					continue
				}

				vs.Logger.Info("Unique count in the last minute",
					"count", count,
					"timestamp", time.Now().Format(time.RFC3339))

			case <-ctx.Done():
				ticker.Stop()
				done <- true
				return
			}
		}
	}()

	<-done
}