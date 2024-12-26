package repository

import (
	"Verve/internal/database"
	"Verve/internal/model/entity"
	"context"
	"fmt"
	"time"
)

const SAVE_ID_KEY_PREFIX = "id"

type VerveRepository interface {
	Save(ctx context.Context, entity entity.VerveEntity) error
	GetUniqueCount(ctx context.Context) (int64, error)
}

type implVerveRepository struct {
	db database.Service
}

func NewImplVerveRepository(database database.Service) *implVerveRepository {
	return &implVerveRepository{
		db: database,
	}
}

func (repo *implVerveRepository) Save(ctx context.Context, entity entity.VerveEntity) error {
	key := fmt.Sprintf("%s%s", SAVE_ID_KEY_PREFIX, entity.Id)
	err := repo.db.Set(ctx, key, entity.Id, 60*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (repo *implVerveRepository) GetUniqueCount(ctx context.Context) (int64, error) {
	count, err := repo.db.CountByPrefix(ctx, SAVE_ID_KEY_PREFIX)
	if err != nil {
		return 0, err
	}
	return count, nil
}
