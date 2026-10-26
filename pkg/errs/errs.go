package errs

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

// Max stack trace depth
const MAX_STACK_TRACE_DEPTH = 5

// ErrorCode should implement this interface
type ErrorCode interface {
	Code() int
	String() string
	Recoverable() bool
}

// CustomError ce
type CustomError struct {
	Code       ErrorCode
	Status     string
	Message    string
	Data       map[string]string
	Additional map[string]any // additional data to be included in the response
	HTTPCode   int            `json:"-"`
	Internal   error          `json:"-"` // Stores the error returned by an external dependency
	stack      []uintptr      // Stack trace for debugging
}

// Error makes it compatible with `error` interface.
func (e *CustomError) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("code=%d, message=%s", e.Code.Code(), e.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%+v", e.Code.Code(), e.Message, e.Internal)
}

// ToHTTPError ToHTTPError
func (e *CustomError) ToHTTPError() *echo.HTTPError {
	var err *echo.HTTPError
	if e.HTTPCode == 0 {
		err = echo.NewHTTPError(http.StatusOK)
	} else {
		err = echo.NewHTTPError(e.HTTPCode)
	}

	msg := map[string]any{}
	msg["code"] = e.Code.Code()
	msg["status"] = e.Status

	if e.Data != nil {
		msg["data"] = e.Data
	}
	if e.Message != "" {
		msg["message"] = e.Message
	}
	if e.Internal != nil {
		err.SetInternal(e.Internal)
		msg["internal"] = e.Internal.Error()
	}

	if e.Additional != nil {
		for k, v := range e.Additional {
			msg[k] = v
		}
	}

	err.Message = msg

	return err
}

// WithHTTPCode .
func (e *CustomError) WithHTTPCode(code int) *CustomError {
	e.HTTPCode = code
	return e
}

// WithMessage WithMessage
func (e *CustomError) WithMessage(message ...string) *CustomError {
	e.Message = strings.Join(message, "")
	return e
}

// WithMessageF WithMessageF
func (e *CustomError) WithMessageF(format string, args ...any) *CustomError {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// WithError WithError
func (e *CustomError) WithError(err error) *CustomError {
	e.Internal = err
	return e
}

// WithoutInternalError .
func (e *CustomError) WithoutInternalError() *CustomError {
	e.Internal = nil
	return e
}

// WithData WithData
func (e *CustomError) WithData(data map[string]string) *CustomError {
	if e.Data == nil {
		e.Data = data
	} else {

		for key, value := range data {
			e.Data[key] = value
		}
	}
	return e
}

// WithDataKV WithData
func (e *CustomError) WithDataKV(key string, value string) *CustomError {
	if e.Data == nil {
		e.Data = make(map[string]string)
	}

	e.Data[key] = value
	return e
}

// ClearData clears all data
func (e *CustomError) ClearData() *CustomError {
	e.Data = nil

	return e
}

// WithAdditional Additional
func (e *CustomError) WithAdditional(key string, value any) *CustomError {
	if e.Additional == nil {
		e.Additional = make(map[string]any)
	}

	e.Additional[key] = value

	return e
}

// WithStack WithData
func (e *CustomError) WithStack() *CustomError {
	return e.WithError(pkgErrors.WithStack(e))
}

// StackTrace returns the stack trace of the error.
func (e *CustomError) StackTrace() string {
	if e.stack == nil {
		return e.Error()
	}

	frames := runtime.CallersFrames(e.stack)
	trace := strings.Builder{}
	trace.WriteString(e.Error())
	trace.WriteString("\n")
	for {
		frame, more := frames.Next()
		// trace.WriteString(fmt.Sprintf("%s %s:%d\n", frame.Function, frame.File, frame.Line))
		trace.WriteString(fmt.Sprintf("%s:%d\n", frame.File, frame.Line))
		if !more {
			break
		}
	}
	return trace.String()
}

// WithStack .
func WithStack(err error, msg ...any) *CustomError {
	if err == nil {
		return nil
	}

	return NewInternalServerError(msg...).WithError(pkgErrors.WithStack(err))
}

// WithStackf .
func WithStackf(err error, format string, args ...any) *CustomError {
	if err == nil {
		return nil
	}

	return WithStack(err, fmt.Sprintf(format, args...))
}

// Is error
func Is(err error, target error) bool {
	a, aok := err.(*CustomError)
	b, bok := target.(*CustomError)
	if aok && bok {
		return a.Code == b.Code
	}

	return errors.Is(err, target)
}

func GetCode(err error) ErrorCode {
	if err == nil {
		//无错误信息
		return StatusOK
	}
	var e *CustomError
	ok := errors.As(err, &e)
	if !ok {
		//非自定义错误格式消息
		return StatusUnknown
	}
	return e.Code
}

// IsCode .
func IsCode(err error, code ErrorCode) bool {
	var custom *CustomError
	if errors.As(err, &custom) && custom != nil {
		return custom.Code.Code() == code.Code()
	}
	return false
}

// IsInstanceNotFound error
func IsInstanceNotFound(err error) bool {
	return IsCode(err, StatusInstanceNotFound)
}

// IsInternalError .
func IsInternalError(err error) bool {
	_, ok := err.(*CustomError)
	if !ok {
		return true
	}

	return IsCode(err, StatusInternalServerError)
}

func IsRecoverable(err error) bool {
	e, ok := err.(*CustomError)
	if ok {
		// recoverable custom error code
		return e.Code.Recoverable()
	}

	return true
}

// NewCustomError NewCustomError
func NewCustomError(code ErrorCode, msg ...any) *CustomError {
	stack := make([]uintptr, MAX_STACK_TRACE_DEPTH)
	length := runtime.Callers(3, stack[:])

	e := &CustomError{
		Code:   code,
		Status: code.String(),
		stack:  stack[:length],
	}

	e.Message = fmt.Sprint(msg...)

	return e
}

// newCustomErrorWithStack NewCustomError
func newCustomErrorWithStack(stack []uintptr, code ErrorCode, msg ...any) *CustomError {
	e := &CustomError{
		Code:   code,
		Status: code.String(),
		stack:  stack,
	}

	e.Message = fmt.Sprint(msg...)

	return e
}

func getStack() []uintptr {
	stack := make([]uintptr, MAX_STACK_TRACE_DEPTH)
	length := runtime.Callers(3, stack[:])
	return stack[:length]
}

func GetStackTrace() string {
	stack := make([]uintptr, MAX_STACK_TRACE_DEPTH)
	length := runtime.Callers(3, stack[:])
	stack = stack[:length]
	frames := runtime.CallersFrames(stack)
	trace := strings.Builder{}
	trace.WriteString("\n")
	for {
		frame, more := frames.Next()
		trace.WriteString(fmt.Sprintf("%s:%d\n", frame.File, frame.Line))
		if !more {
			break
		}
	}
	return trace.String()
}

// NewCustomErrorWithoutStack NewCustomError
func NewCustomErrorWithoutStack(code ErrorCode, msg ...any) *CustomError {
	e := &CustomError{
		Code:   code,
		Status: code.String(),
	}

	e.Message = fmt.Sprint(msg...)

	return e
}

// NewCustomErrorf NewCustomError
func NewCustomErrorf(code ErrorCode, format string, msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), code, fmt.Sprintf(format, msg...))
}

// NewCustomErrorf NewCustomError
func newCustomErrorWithStackf(stack []uintptr, code ErrorCode, format string, msg ...any) *CustomError {
	return newCustomErrorWithStack(stack, code, fmt.Sprintf(format, msg...))
}

// NewBadRequestError .
func NewBadRequestError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusBadRequest, msg...)
}

func NewBadRequestErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusBadRequest, format, msg...)
}

// NewUnauthorizedError .
func NewUnauthorizedError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusUnauthorized, msg...)
}

// NewNotFoundError .
func NewNotFoundError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusNotFound, msg...)
}

// NewInternalServerError .
func NewInternalServerError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusInternalServerError, msg...).WithHTTPCode(http.StatusInternalServerError)
}

// NewInternalServerErrorf .
func NewInternalServerErrorf(format string, msg ...any) *CustomError {
	return NewCustomErrorf(StatusInternalServerError, format, msg...).WithHTTPCode(http.StatusInternalServerError)
}

// NewForbiddenError .
func NewForbiddenError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusForbidden, msg...)
}

// NewServerMaintainError .
func NewServerMaintainError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusServerMaintain, msg...)
}

// NewInstanceNotFoundError .
func NewInstanceNotFoundError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusInstanceNotFound, msg...)
}

// NewInstanceNotFoundError .
func NewInstanceNotFoundErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusInstanceNotFound, format, msg...)
}

// NewNotImplementError .
// NewNotImplementError .
func NewNotImplementError() *CustomError {
	return newCustomErrorWithStack(getStack(), StatusInternalServerError, "method not implement")
}

// NewNotEnoughResourceError .
func NewNotEnoughResourceError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusNotEnoughResource, msg...)
}

// NewDuplicatedError .
func NewDuplicatedError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusDuplicated, msg...)
}

// 自定义error 信息
func MyErrors(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}
