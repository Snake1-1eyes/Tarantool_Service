package tarantool

import (
	// ...existing imports...
	"fmt"

	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/models"
	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/utils"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
)

type TarantoolRepo struct {
	conn *tarantool.Connection
	log  *zap.Logger
}

func NewTarantoolRepo(conn *tarantool.Connection, log *zap.Logger) *TarantoolRepo {
	return &TarantoolRepo{
		conn: conn,
		log:  log,
	}
}

func (r *TarantoolRepo) Create(kv *models.KeyValue) (*models.KeyValue, error) {
	r.log.Info("Создание новой пары ключ-значение",
		zap.String("key", kv.Key),
		zap.Any("value", kv.Value),
	)

	req := tarantool.NewInsertRequest("kv").Tuple([]interface{}{kv.Key, kv.Value})

	future := r.conn.Do(req)
	resp, err := future.Get()

	if err != nil {
		if tarErr, ok := err.(*tarantool.Error); ok && tarErr.Code == 3 {
			r.log.Warn("Попытка создать существующий ключ",
				zap.String("key", kv.Key),
				zap.Error(err),
			)
			return nil, fmt.Errorf("key already exists")
		}

		r.log.Error("Ошибка создания записи",
			zap.String("key", kv.Key),
			zap.Error(err),
		)
		return nil, fmt.Errorf("key already exists")
	}

	tuple := resp[0].([]interface{})
	result := &models.KeyValue{
		Key:   tuple[0].(string),
		Value: utils.ConvertMapToStringKeys(tuple[1]),
	}

	r.log.Info("Запись успешно создана",
		zap.String("key", result.Key),
		zap.Any("value", result.Value),
	)
	return result, nil
}

func (r *TarantoolRepo) Get(key string) (*models.KeyValue, error) {
	r.log.Info("Получение значения по ключу",
		zap.String("key", key),
	)

	req := tarantool.NewSelectRequest("kv").
		Index("primary").
		Key([]interface{}{key})

	future := r.conn.Do(req)
	resp, err := future.Get()
	if err != nil {
		r.log.Error("Ошибка получения значения",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, err
	}

	if len(resp) == 0 {
		r.log.Info("Ключ не найден",
			zap.String("key", key),
		)
		return nil, fmt.Errorf("key not found")
	}

	tuple := resp[0].([]interface{})
	value := utils.ConvertMapToStringKeys(tuple[1])

	result := &models.KeyValue{
		Key:   tuple[0].(string),
		Value: value,
	}

	r.log.Info("Значение успешно получено",
		zap.String("key", result.Key),
		zap.Any("value", result.Value),
	)
	return result, nil
}

func (r *TarantoolRepo) Update(key string, value interface{}) (*models.KeyValue, error) {
	r.log.Info("Обновление значения",
		zap.String("key", key),
		zap.Any("new_value", value),
	)

	convertedValue := utils.ConvertMapToStringKeys(value)

	req := tarantool.NewUpdateRequest("kv").
		Index("primary").
		Key([]interface{}{key}).
		Operations(tarantool.NewOperations().Assign(1, convertedValue))

	future := r.conn.Do(req)
	resp, err := future.Get()
	if err != nil {
		r.log.Error("Ошибка обновления значения",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, err
	}

	if len(resp) == 0 {
		r.log.Info("Ключ не найден при обновлении",
			zap.String("key", key),
		)
		return nil, fmt.Errorf("key not found")
	}

	tuple := resp[0].([]interface{})
	result := &models.KeyValue{
		Key:   tuple[0].(string),
		Value: utils.ConvertMapToStringKeys(tuple[1]),
	}

	r.log.Info("Значение успешно обновлено",
		zap.String("key", result.Key),
		zap.Any("value", result.Value),
	)
	return result, nil
}

func (r *TarantoolRepo) Delete(key string) error {
	r.log.Info("Удаление записи",
		zap.String("key", key),
	)

	_, err := r.Get(key)
	if err != nil {
		return err
	}

	req := tarantool.NewDeleteRequest("kv").
		Index("primary").
		Key([]interface{}{key})

	future := r.conn.Do(req)
	resp, err := future.Get()
	if err != nil {
		r.log.Error("Ошибка удаления записи",
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	if len(resp) == 0 {
		r.log.Info("Ключ не найден при удалении",
			zap.String("key", key),
		)
		return fmt.Errorf("key not found")
	}

	r.log.Info("Запись успешно удалена",
		zap.String("key", key),
	)
	return nil
}
