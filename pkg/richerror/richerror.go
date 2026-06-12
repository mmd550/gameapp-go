package richerror

type ErrorKind uint

const (
	KindInvalid ErrorKind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type RichError struct {
	wrappedError error
	message      string
	kind         ErrorKind
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

func (err RichError) WithKind(kind ErrorKind) RichError {
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

func (err RichError) Kind() ErrorKind {
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
	return r.message
}
