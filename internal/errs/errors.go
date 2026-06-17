package errs

import "fmt"

type Error int

const (
	InternalError                   Error = -1
	NegativeSubscribePriceError     Error = 0
	WrongSubscribeDateIntervalError Error = 1
	NotFoundSubscribeError          Error = 2
	JsonValidationError             Error = 3
)

func (e Error) Error() string {
	title := e.Title()
	if title == "" {
		return fmt.Sprintf("unknown error code %d", e)
	}
	return title
}

func (e Error) Title() string {
	switch e {
	case InternalError:
		return "internal error"
	case NegativeSubscribePriceError:
		return "negative subscribe price"
	case WrongSubscribeDateIntervalError:
		return "wrong subscribe date interval"
	case NotFoundSubscribeError:
		return "subscribe not found"
	case JsonValidationError:
		return "json validation error"
	}
	return ""
}
