package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Email", err.Field()))
		case "password":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Password", err.Field()))
		case "name":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Name", err.Field()))
		case "description":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Description", err.Field()))
		case "user_id":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid UserID", err.Field()))
		case "public":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Public", err.Field()))
		case "channel_id":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid ChannelID", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
