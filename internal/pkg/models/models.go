package models

type KeyValue struct {
	Key   string      `json:"key" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

type UpdateValue struct {
	Value interface{} `json:"value" validate:"required"`
}
