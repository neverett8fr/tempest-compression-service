package application

type Response struct {
	Errors []error     `json:"errors"`
	Data   interface{} `json:"data"`
}

func NewResponse(data interface{}, err ...error) *Response {

	return &Response{
		Errors: err,
		Data:   data,
	}
}
