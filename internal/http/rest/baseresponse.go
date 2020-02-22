package rest

type BaseResponse struct {
	Data   interface{}
	Error  []*APIErrors
	Status int
}

type APIErrors struct {
	Detail string `json:"detail"`
}

func NewBaseResponse(status int, data interface{}, error interface{}) *BaseResponse {
	errs := ParseErrors(error)
	return &BaseResponse{
		Data:   data,
		Error:  errs,
		Status: status,
	}
}

func ParseErrors(errs interface{}) []*APIErrors {

	errorList, _ := errs.([]error)
	var errors []string
	for x := 0; x < len(errorList); x++ {
		errors = append(errors, errorList[x].Error())
	}

	var parsedErrorList []*APIErrors

	for i := 0; i < len(errors); i++ {
		tempErr := APIErrors{Detail:errors[i]}
		parsedErrorList = append(parsedErrorList, &tempErr)
	}

	return parsedErrorList
}


