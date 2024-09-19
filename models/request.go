package models

type RequestID struct {
	ID int64 `json:"id" validate:"required"`
}

type TestingHandlerRequest struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Body   interface{} `json:"body"`
}

type TestingHandlerExpected struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}

