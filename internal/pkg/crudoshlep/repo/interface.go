package repo

import "github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/models"

type Repository interface {
	Create(kv *models.KeyValue) (*models.KeyValue, error)
	Get(key string) (*models.KeyValue, error)
	Update(key string, value any) (*models.KeyValue, error)
	Delete(key string) error
}
