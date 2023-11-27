package viewmodel

type TestBodyViewModelRequest struct {
	// in: body
	Body struct {
		Field string `json:"field" binding:"required,min=1"`
	}
}

type TestUrlParamsViewModelRequest struct {
	Field string `form:"field" binding:"required,min=1"`
}

type TestViewModelResponse struct {
	// in: body
	Body struct {
		Field string `json:"field"`
	} `json:"body"`
}

type TestNoBodyViewModelResponse struct{}
