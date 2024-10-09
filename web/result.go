package web

import (
	"fmt"
)

type Result struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func ErrorResult(message *MessageCode) *Result {
	return &Result{message.Code, message.Message, nil}
}

func OK(params ...interface{}) *Result {
	size := len(params)
	switch size {
	case 0:
		return &Result{0, "", nil}
	case 1:
		if data, ok := params[0].(map[string]interface{}); ok {
			return &Result{0, "", data}
		}
	case 2:
		if data, ok := params[0].(string); ok {
			return &Result{0, "", map[string]interface{}{
				data: params[1],
			}}
		}
	default:
		i := size % 2
		if i == 0 {
			data := make(map[string]interface{})
			for i = 0; i < size; i += 2 {
				if key, ok := params[i].(string); ok {
					data[key] = params[i+1]
				} else {
					break
				}
			}

			return &Result{0, "", data}
		}
	}
	return &Result{0, "", nil}
}

func Error(params ...interface{}) *Result {
	size := len(params)

	if size > 0 {
		var code int
		var message string
		if data, ok := params[0].(int); ok {
			code = data
		} else if data, ok := params[0].(*MessageCode); ok {
			code = data.Code
			message = data.Message
		}

		if size > 1 {
			if message == "" {
				if data, ok := params[1].(string); ok {
					message = data
				}
			} else {
				message = fmt.Sprintf(message, params[1:]...)
			}
		}

		return &Result{code, message, nil}
	}

	return ErrorResult(ERROR_SYSTEM)
}
