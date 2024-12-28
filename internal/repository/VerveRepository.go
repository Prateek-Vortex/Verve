package repository

import (
	"Verve/internal/database"
	"Verve/internal/model/entity"
	"context"
)

const SAVE_ID_KEY = "id"

type VerveRepository interface {
	Save(ctx context.Context, entity entity.VerveEntity) error
	GetUniqueCount(ctx context.Context) (int64, error)
	Delete(ctx context.Context) error
}

type implVerveRepository struct {
	db database.Service
}

func NewImplVerveRepository(database database.Service) *implVerveRepository {
	return &implVerveRepository{
		db: database,
	}
}

// Uncomment below code to implement the unique count based on the id of the entity according to the sliding time window size.
// Eample if key is "Prateek"  which comes  multiple times in 1-59sec window then the count will be 1, and id "Prateek" also comes
// at the 59th sec than the next minute window size will also be one till 1-58sec of the next minute.

// func (repo *implVerveRepository) Save(ctx context.Context, entity entity.VerveEntity) error {
// 	key := fmt.Sprintf("%s%s", SAVE_ID_KEY, entity.Id)
// 	err := repo.db.Set(ctx, key, entity.Id, 60*time.Second)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (repo *implVerveRepository) GetUniqueCount(ctx context.Context) (int64, error) {
// 	count, err := repo.db.CountByPrefix(ctx, SAVE_ID_KEY)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// This is the implementation of the unique count based on the id of the entity according to the HyperLogLog algorithm,
// which is a probabilistic data structure used to count unique elements in a set. This will count the id base on
// fixed window size that is 1-59 sec.

func (repo *implVerveRepository) Save(ctx context.Context, entity entity.VerveEntity) error {

	err := repo.db.SAdd(ctx, SAVE_ID_KEY, entity.Id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *implVerveRepository) GetUniqueCount(ctx context.Context) (int64, error) {
	return repo.db.SCard(ctx, SAVE_ID_KEY)
}

func (repo *implVerveRepository) Delete(ctx context.Context) error {
	return repo.db.Del(ctx, SAVE_ID_KEY)
}
