package error

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"net/http"
)

type HttpErrorResponse struct {
	Typ string `json:"type"`
	Msg string `json:"msg"`
}

type HttpError struct {
	err    error
	status int
	res    HttpErrorResponse
}

type Error struct {
	err  error
	code Code
	msg  string
}

func NewHttpError(err error) HttpError {

	var e HttpError
	e.status = http.StatusInternalServerError
	e.err = errors.New("internal error")
	e.res = makeErrorResponse(CodeUnknown, "internal error")

	fmt.Printf("%w", err)
	var ee *Error
	if errors.As(UnwrapErrChain(err), &ee) {
		switch ee.Code() {
		case CodeNotFound:
			e.status = http.StatusNotFound
		case CodeInvalidArgument,
			CodeDbDuplication,
			CodeValidation,
			CodeDbReference:
			e.status = http.StatusBadRequest
		case CodeJwtMalformed,
			CodeJwtInvalid:
			e.status = http.StatusUnauthorized
		case CodeUnknown:
			fallthrough
		default:
			e.status = http.StatusInternalServerError
		}

		e.res = makeErrorResponse(ee.Code(), ee.Error())
	}

	return e
}

func (h *HttpError) JSON() HttpErrorResponse {
	return h.res
}

func (h *HttpError) Status() int {
	return h.status
}

func makeErrorResponse(code Code, msg string) HttpErrorResponse {
	return HttpErrorResponse{
		Typ: fmt.Sprintf("Error.%s", code),
		Msg: msg,
	}
}

func NewError(code Code, msg string) error {
	return WrapError(nil, code, msg)
}

func WrapError(err error, code Code, msg string) error {
	return &Error{
		err:  err,
		code: code,
		msg:  msg,
	}
}

func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s", e.msg)
	}

	return e.msg
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Code() Code {
	return e.code
}

func HandlePgError(err error) error {
	e := NewError(CodeUnknown, err.Error())

	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return NewError(CodeDbDuplication, pgErr.Detail)
		case "23502":
			return NewError(CodeDbReference, pgErr.Message)
		default:
			return e
		}
	}

	return e
}

func UnwrapErrChain(err error) error {
	var t *Error
	for err != nil {
		if errors.As(err, &t) {
			err = t.Unwrap()
		} else {
			err = nil
		}
	}
	return t
}
