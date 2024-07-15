package repository

import (
	"fmt"
	"strings"
)

type UniqueConstrainError struct {
	constrain string
	columns   []string
	err       error
}

func (e *UniqueConstrainError) Unwrap() error {
	return e.err
}

func (e *UniqueConstrainError) Error() string {
	return fmt.Sprintf("unique constrain %s failed", e.constrain)
}

func (e *UniqueConstrainError) FailedConstrain() string {
	return e.constrain
}

func (e *UniqueConstrainError) Columns() string {
	return strings.Join(e.columns, ",")
}

func newUniqueConstrainError(err error, constrain string, columns ...string) *UniqueConstrainError {
	return &UniqueConstrainError{
		err:       err,
		constrain: constrain,
		columns:   columns,
	}
}
