package main

import (
	"fmt"
)

type (
	// ErrCode type
	ErrCode uint16
	// appErr struct
	appErr struct {
		Code     ErrCode `json:"code,omitempty"`
		HTTPCode int     `json:"httpCode,omitempty"`
		InnerErr error   `json:"innerErr,omitempty"`
	}
)

// Error codes
const (
	UnknownErrCode ErrCode = iota

	ImageSizeTooBigErrCode
	ImageAlreadyExistsErrCode
	InvalidBucketNameErrCode
	InvalidImageErrCode
)

// Error returns message by error code
func (e *appErr) Error() string {
	return ""
}

// Stack returns error's stack trace
func (e *appErr) Stack() string {
	return fmt.Sprintf("%+v", e.InnerErr)
}

// func isNotFoundErr(err error) bool {
// 	if e, ok := errors.Cause(err).(*Error); ok && e.Code == NotFoundErrCode {
// 		return true
// 	}
// 	return false
// }
