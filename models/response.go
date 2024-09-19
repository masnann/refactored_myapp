package models

import "time"

type Response struct {
	StatusCode       string      `json:"statusCode"`
	Success          bool        `json:"success"`
	ResponseDateTime time.Time   `json:"responseDateTime"`
	Result           interface{} `json:"result"`
	Message          string      `json:"message"`
}

type ResponseLogError struct {
	HttpCode      int    `json:"httpCode"`
	StatusCode    string `json:"errorCode"`
	UserMessage   string `json:"userMessage"`
	SystemMessage string `json:"systemMessage"`
}
