package service

import (
	restclient "Verve/internal/configs/restClient"
	"Verve/internal/event"
	"Verve/internal/model/entity"
	"Verve/internal/model/request"
	"Verve/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

type VerveService interface {
	SaveAndPost(ctx context.Context, verveRequest request.VerveRequest) error
	LogUniqueCountEveryMinute(ctx context.Context)
	SendUniqueCountEveryMinute(ctx context.Context)
}

type implVerveService struct {
	verveRepo  repository.VerveRepository
	restClient restclient.RestClient
	Logger     *slog.Logger
	Event      event.Event
}

func NewImplVerveService(repository repository.VerveRepository, client restclient.RestClient, logger *slog.Logger, event event.Event) *implVerveService {
	return &implVerveService{
		verveRepo:  repository,
		restClient: client,
		Logger:     logger,
		Event:      event,
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
		vs.Logger.Error("Failed to get count", "error", err.Error())
	}
	jsonData := map[string]int64{"count": count}
	_, err = vs.restClient.Post(url, jsonData)
	if err != nil {
		vs.Logger.Error("Failed to send post request to url", url, err)
	}
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

func (vs *implVerveService) SendUniqueCountEveryMinute(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				vs.Logger.Error("Recovered from panic", "error", r)
			}
		}()
		for {
			select {
			case <-ticker.C:
				count, err := vs.verveRepo.GetUniqueCount(ctx)
				if err != nil {
					vs.Logger.Error("Failed to get unique count", "error", err)
					continue
				}
				err = vs.verveRepo.Delete(ctx)
				if err != nil {
					vs.Logger.Error("Failed to delete unique count", "error", err)
					errs := fmt.Errorf("failed to delete unique count %w", err)
					panic(errs)
				}

				vs.Event.Publish(ctx, "unique_count", strconv.FormatInt(count, 10))

			case <-ctx.Done():
				ticker.Stop()
				done <- true
				return
			}
		}
	}()

	<-done
}
