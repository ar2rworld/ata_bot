package myerror

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MyError struct {
	message string
}

func (e *MyError) Error () string {
	return e.message
}

func NewError(m string) *MyError {
	return &MyError{message: m,}
}

func NewAPIResponseError(res *tgbotapi.APIResponse, e error) error {
	return NewError(
		fmt.Sprintf(`APIResponse error: "%s", errorCode: %d, description: %s, json res: %s`,
		e, res.ErrorCode, res.Description, res.Result),
	)
}
