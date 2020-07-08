package errors

type GoLearningError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e GoLearningError) Error() string {
	return e.Message
}
