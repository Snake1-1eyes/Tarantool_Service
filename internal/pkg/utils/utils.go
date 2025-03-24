package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Snake1-1eyes/Tarantool_Service/internal/pkg/models"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func GetRequestData(r *http.Request, requestData any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestData); err != nil {
		return err
	}

	return nil
}

func WriteResponseData(w http.ResponseWriter, responseData interface{}, successStatusCode int) error {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(responseData)
	if err != nil {
		return fmt.Errorf("error marshalling response: %w", err)
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(successStatusCode)

	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("error writing response: %w", err)
	}

	return nil
}

func WriteErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, `{"message":"%s"}`, message)
}

func ConvertMapToStringKeys(m interface{}) interface{} {
	switch val := m.(type) {
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range val {
			strKey := fmt.Sprint(k)
			newMap[strKey] = ConvertMapToStringKeys(v)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(val))
		for i, v := range val {
			newSlice[i] = ConvertMapToStringKeys(v)
		}
		return newSlice
	default:
		return val
	}
}

func ValidateKeyValue(kv *models.KeyValue) error {
	if kv.Key == "" {
		return fmt.Errorf("ключ не может быть пустым")
	}
	if kv.Value == nil {
		return fmt.Errorf("значение не может быть пустым")
	}

	switch v := kv.Value.(type) {
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("значение не может быть пустым объектом")
		}
	case map[interface{}]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("значение не может быть пустым объектом")
		}
	default:
		return fmt.Errorf("значение должно быть JSON объектом")
	}

	return nil
}

func ValidateUpdateValue(update *models.UpdateValue) error {
	if update.Value == nil {
		return fmt.Errorf("значение не может быть пустым")
	}

	switch v := update.Value.(type) {
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("значение не может быть пустым объектом")
		}
	case map[interface{}]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("значение не может быть пустым объектом")
		}
	default:
		return fmt.Errorf("значение должно быть JSON объектом")
	}

	return nil
}
