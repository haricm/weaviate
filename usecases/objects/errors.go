//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package objects

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	// ErrItemNotFound item doesn't exist
	ErrItemNotFound = errors.New("item not found")
	// ErrBadRequest because of invalid inputs
	ErrBadRequest = errors.New("validation")
	//  ErrAccessDenied access denied
	ErrAccessDenied = errors.New("authorization")
	// ErrServiceInternal is an internal service error
	ErrServiceInternal = errors.New("service internal")
)

func IsErrorNotFound(err error) bool {
	return errors.Is(err, ErrItemNotFound)
}

func IsErrorAccessDenied(err error) bool {
	return errors.Is(err, ErrAccessDenied)
}

func IsErrorBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

func IsErrorInternal(err error) bool {
	return errors.Is(err, ErrServiceInternal)
}

// ErrInvalidUserInput indicates a client-side error
type ErrInvalidUserInput struct {
	msg string
}

func (e ErrInvalidUserInput) Error() string {
	return e.msg
}

// NewErrInvalidUserInput with Errorf signature
func NewErrInvalidUserInput(format string, args ...interface{}) ErrInvalidUserInput {
	return ErrInvalidUserInput{msg: fmt.Sprintf(format, args...)}
}

// ErrInternal indicates something went wrong during processing
type ErrInternal struct {
	msg string
}

func (e ErrInternal) Error() string {
	return e.msg
}

// NewErrInternal with Errorf signature
func NewErrInternal(format string, args ...interface{}) ErrInternal {
	return ErrInternal{msg: fmt.Sprintf(format, args...)}
}

// ErrNotFound indicates the desired resource doesn't exist
type ErrNotFound struct {
	msg string
}

func (e ErrNotFound) Error() string {
	return e.msg
}

// NewErrNotFound with Errorf signature
func NewErrNotFound(format string, args ...interface{}) ErrNotFound {
	return ErrNotFound{msg: fmt.Sprintf(format, args...)}
}
