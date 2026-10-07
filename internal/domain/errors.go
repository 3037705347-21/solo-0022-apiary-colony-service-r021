package domain

import (
	"errors"
	"fmt"
)

// ErrorCode 是可以稳定映射到 HTTP 状态码的领域错误分类。
type ErrorCode string

const (
	CodeInvalidRequest ErrorCode = "invalid_request"
	CodeNotFound       ErrorCode = "not_found"
	CodeConflict       ErrorCode = "conflict"
	CodeInvalidState   ErrorCode = "invalid_state"
	CodeStorage        ErrorCode = "storage_error"
)

// Error 携带稳定分类、可选字段名与面向调用方的说明。
type Error struct {
	Code    ErrorCode
	Field   string
	Message string
}

func (e *Error) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func Invalid(field, format string, args ...any) *Error {
	return &Error{Code: CodeInvalidRequest, Field: field, Message: fmt.Sprintf(format, args...)}
}

func NotFound(kind, id string) *Error {
	return &Error{Code: CodeNotFound, Message: fmt.Sprintf("%s %q 不存在", kind, id)}
}

func Conflict(format string, args ...any) *Error {
	return &Error{Code: CodeConflict, Message: fmt.Sprintf(format, args...)}
}

func InvalidState(format string, args ...any) *Error {
	return &Error{Code: CodeInvalidState, Message: fmt.Sprintf(format, args...)}
}

func Storage(format string, args ...any) *Error {
	return &Error{Code: CodeStorage, Message: fmt.Sprintf(format, args...)}
}

// AsError 提取链路上的领域错误，供入口层映射响应。
func AsError(err error) (*Error, bool) {
	var target *Error
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

func CodeOf(err error) ErrorCode {
	if target, ok := AsError(err); ok {
		return target.Code
	}
	return CodeStorage
}
