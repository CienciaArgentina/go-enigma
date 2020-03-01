package rest

type BaseResponse struct {
	Data   interface{}  `json:"data"`
	Error  []*APIErrors `json:"error"`
	Status int          `json:"status"`
}

type APIErrors struct {
	Detail string `json:"detail"`
}

func NewBaseResponse(status int, data interface{}, error interface{}) *BaseResponse {
	var errs []*APIErrors
	if error != nil {
		errs = ParseErrors(error)
	}

	return &BaseResponse{
		Data:   data,
		Error:  errs,
		Status: status,
	}
}

func ParseErrors(errs interface{}) []*APIErrors {
	var parsedErrorList []*APIErrors
	errorList, isErrorList := errs.([]error)
	if isErrorList {
		var errors []string
		for x := 0; x < len(errorList); x++ {
			errors = append(errors, errorList[x].Error())
		}

		for i := 0; i < len(errors); i++ {
			tempErr := APIErrors{Detail: errors[i]}
			parsedErrorList = append(parsedErrorList, &tempErr)
		}
	} else {
		err := errs.(error)
		parsedErrorList = append(parsedErrorList, &APIErrors{Detail: err.Error()})
	}

	return parsedErrorList
}
