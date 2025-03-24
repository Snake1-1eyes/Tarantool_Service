package usecase

import (
	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/crudoshlep/repo"
	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/models"
)

type KVUseCase struct {
	repo repo.Repository
}

func NewKVUseCase(repo repo.Repository) *KVUseCase {
	return &KVUseCase{repo: repo}
}

func (uc *KVUseCase) Create(kv *models.KeyValue) (*models.KeyValue, error) {
	return uc.repo.Create(kv)
}

func (uc *KVUseCase) Get(key string) (*models.KeyValue, error) {
	return uc.repo.Get(key)
}

func (uc *KVUseCase) Update(key string, value any) (*models.KeyValue, error) {
	return uc.repo.Update(key, value)
}

func (uc *KVUseCase) Delete(key string) error {
	return uc.repo.Delete(key)
}
