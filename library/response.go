package library

type AppResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResp(msg string, data interface{}) AppResponse {

	return AppResponse{
		Message: msg,
		Data:    data,
	}
}
