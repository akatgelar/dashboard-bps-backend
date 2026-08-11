package models

type BaseResponse struct {
	Status		bool `json:"status"`
	Message		string `json:"message"`
	Data		any `json:"data"`
	Metadata	any `json:"metadata"`
}