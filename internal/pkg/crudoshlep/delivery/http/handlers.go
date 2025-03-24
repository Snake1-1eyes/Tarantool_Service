package http

import (
	"net/http"

	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/crudoshlep/usecase"
	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/models"
	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	useCase *usecase.KVUseCase
	log     *zap.Logger
}

func NewHandler(useCase *usecase.KVUseCase, log *zap.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		log:     log,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var kv models.KeyValue
	if err := utils.GetRequestData(r, &kv); err != nil {
		h.log.Error("Ошибка разбора тела запроса",
			zap.Error(err),
			zap.String("method", "Create"),
		)
		utils.WriteErrorMessage(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	if err := utils.ValidateKeyValue(&kv); err != nil {
		h.log.Warn("Ошибка валидации",
			zap.Error(err),
			zap.String("key", kv.Key),
			zap.Any("value", kv.Value),
		)
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.useCase.Create(&kv)
	if err != nil {
		if err.Error() == "key already exists" {
			h.log.Warn("Попытка создать существующий ключ",
				zap.String("key", kv.Key),
			)
			utils.WriteErrorMessage(w, http.StatusConflict, "ключ уже существует")
			return
		}
		h.log.Error("Ошибка создания записи",
			zap.Error(err),
			zap.String("key", kv.Key),
		)
		utils.WriteErrorMessage(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	h.log.Info("Запись успешно создана",
		zap.String("key", result.Key),
		zap.Any("value", result.Value),
	)
	utils.WriteResponseData(w, result, http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	h.log.Info("Получен запрос на чтение",
		zap.String("key", key),
	)

	kv, err := h.useCase.Get(key)
	if err != nil {
		if err.Error() == "key not found" {
			h.log.Info("Ключ не найден",
				zap.String("key", key),
			)
			utils.WriteErrorMessage(w, http.StatusNotFound, "ключ не найден")
			return
		}
		h.log.Error("Ошибка получения записи",
			zap.Error(err),
			zap.String("key", key),
		)
		utils.WriteErrorMessage(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	h.log.Info("Запись успешно получена",
		zap.String("key", kv.Key),
		zap.Any("value", kv.Value),
	)
	utils.WriteResponseData(w, kv, http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var update models.UpdateValue
	if err := utils.GetRequestData(r, &update); err != nil {
		h.log.Error("Ошибка разбора тела запроса",
			zap.Error(err),
			zap.String("key", key),
		)
		utils.WriteErrorMessage(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	if err := utils.ValidateUpdateValue(&update); err != nil {
		h.log.Warn("Ошибка валидации",
			zap.Error(err),
			zap.String("key", key),
			zap.Any("value", update.Value),
		)
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.useCase.Update(key, update.Value)
	if err != nil {
		if err.Error() == "key not found" {
			h.log.Info("Ключ не найден при обновлении",
				zap.String("key", key),
			)
			utils.WriteErrorMessage(w, http.StatusNotFound, "ключ не найден")
			return
		}
		h.log.Error("Ошибка обновления записи",
			zap.Error(err),
			zap.String("key", key),
		)
		utils.WriteErrorMessage(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	h.log.Info("Запись успешно обновлена",
		zap.String("key", result.Key),
		zap.Any("value", result.Value),
	)
	utils.WriteResponseData(w, result, http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	h.log.Info("Получен запрос на удаление",
		zap.String("key", key),
	)

	err := h.useCase.Delete(key)
	if err != nil {
		if err.Error() == "key not found" {
			h.log.Info("Ключ не найден при удалении",
				zap.String("key", key),
			)
			utils.WriteErrorMessage(w, http.StatusNotFound, "ключ не найден")
			return
		}
		h.log.Error("Ошибка удаления записи",
			zap.Error(err),
			zap.String("key", key),
		)
		utils.WriteErrorMessage(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	h.log.Info("Запись успешно удалена",
		zap.String("key", key),
	)
	utils.WriteResponseData(w, nil, http.StatusOK)
}
