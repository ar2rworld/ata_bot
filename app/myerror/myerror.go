package myerror

type MyError struct {
	message string
}

func (e *MyError) Error () string {
	return e.message
}

func NewError(m string) *MyError {
	return &MyError{message: m,}
}