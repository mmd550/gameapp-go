package richerror

import "fmt"

type Kind uint

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type RichError struct {
	wrappedError error
	message      string
	kind         Kind
	operation    string
	meta         map[string]any
}

func New(operation string) RichError {
	return RichError{
		operation: operation,
	}
}

func (err RichError) WithMessage(message string) RichError {
	err.message = message
	return err
}

func (err RichError) WithKind(kind Kind) RichError {
	err.kind = kind
	return err
}

func (err RichError) WithErr(e error) RichError {
	err.wrappedError = e
	return err
}

func (err RichError) WithMeta(meta map[string]any) RichError {
	err.meta = meta
	return err
}

func (err RichError) Kind() Kind {
	if err.kind != 0 {
		return err.kind
	}

	re, ok := err.wrappedError.(RichError)

	if !ok {
		return KindUnexpected
	}

	return re.Kind()
}

func (err RichError) Message() string {
	if err.message != "" {
		return err.message
	}

	re, ok := err.wrappedError.(RichError)

	if !ok {
		return "unexpected error"
	}

	return re.Message()
}

func (r RichError) Error() string {
	return fmt.Sprintf("%s: %s", r.operation, r.message)
}
