package errors

import (
	"fmt"

	"github.com/go-errors/errors"
)

func WithStack(err error) error {
	return errors.Wrap(err, 1)
}

func Wrap(err error, msg string) error {
	return errors.WrapPrefix(err, msg, 1)
}

func Wrapf(err error, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return Wrap(err, msg)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}
