package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

func NewFullErrorResponse(statusCode int, err error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		RootErr:    err,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}
func NewErrorResponse(err error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    err,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}
func NewUnauthorized(err error, msg, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    err,
		Message:    msg,
		Key:        key,
	}
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}

	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

func NewCustomError(root error, msg string, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}

	return NewErrorResponse(errors.New(msg), msg, msg, key)
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong in DB", err.Error(), "DB_ERROR")
}

func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrInvalidObjectId(err error) *AppError {
	return NewErrorResponse(err, "invalid objectId", err.Error(), "ErrInvalidObjectId")
}

func ErrInvalidPostStatus(err error) *AppError {

	return NewErrorResponse(err, "invalid post status", err.Error(), "ErrInvalidPostStatus")
}

func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong in the server", err.Error(), "ErrInternal")
}

func ErrCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotList%s", entity),
	)
}

func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotCreate%s", entity),
	)
}
func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotDelete%s", entity),
	)
}

func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotUpdate%s", entity),
	)
}
func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("Cannot get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGet%s", entity),
	)
}
func ErrEntityDeleted(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("%s deleted", strings.ToLower(entity)),
		fmt.Sprintf("Err%sDeleted", entity),
	)
}
func ErrEntityExists(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("%s already exists", strings.ToLower(entity)),
		fmt.Sprintf("Err%sAlreadyExists", entity),
	)
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("%s not found", strings.ToLower(entity)),
		fmt.Sprintf("Err%sAlreadyNotFound", entity),
	)
}

func ErrNoPermission(err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("You don't have permission: %s", err),
		"ErrNoPermission",
	)
}
func ErrUnauthorized(err error) *AppError {
	return NewCustomError(err,
		fmt.Sprintf("You don't have permission:%s", err),
		"ErrUnauthorized",
	)
}
func ErrExpiredToken() *AppError {
	return NewCustomError(errors.New("token expired"), "token expired", "ErrExpiredToken")
}
func ErrInvalidToken() *AppError {
	return NewCustomError(errors.New("token is invalid"), "token is invalid", "ErrInvalidToken")
}

var RecordNotFound = errors.New("record not found")

func ErrResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}
