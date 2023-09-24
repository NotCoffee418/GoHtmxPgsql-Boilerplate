package common

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

var ApiResponseFactory = new(apiResponseAlias)

type apiResponseAlias ApiResponse

func (a *apiResponseAlias) Ok(data interface{}) *ApiResponse {
	return &ApiResponse{
		Success: true,
		Data:    data,
	}
}

func (a *apiResponseAlias) Fail(err error) *ApiResponse {
	return &ApiResponse{
		Success: false,
		Error:   err,
	}
}
