package dto

type BaseResponse[T any] struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    T      `json:"data"`
}
