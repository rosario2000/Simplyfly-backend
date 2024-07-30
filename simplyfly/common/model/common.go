package model

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	HttpStatus int    `json:"httpStatus"`
}

type KafkaMessage string
