package rest

type BaseResponse struct {
	Data   interface{}
	Error  interface{}
	Status int
}

func NewBaseResponse(status int, data interface{}, error interface{}) *BaseResponse {
	return &BaseResponse{
		Data:   data,
		Error:  error,
		Status: status,
	}
}
