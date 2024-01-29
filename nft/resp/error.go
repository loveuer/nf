package resp

import (
	"errors"
	"fmt"
	"github.com/loveuer/nf"
)

type Error struct {
	status uint32
	msg    string
	err    error
	data   any
}

func (e Error) Error() string {
	if e.msg != "" {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}

	switch e.status {
	case 200:
		return fmt.Sprintf("%s: %s", MSG200, e.err.Error())
	case 202:
		return fmt.Sprintf("%s: %s", MSG202, e.err.Error())
	case 400:
		return fmt.Sprintf("%s: %s", MSG400, e.err.Error())
	case 401:
		return fmt.Sprintf("%s: %s", MSG401, e.err.Error())
	case 403:
		return fmt.Sprintf("%s: %s", MSG403, e.err.Error())
	case 404:
		return fmt.Sprintf("%s: %s", MSG404, e.err.Error())
	case 429:
		return fmt.Sprintf("%s: %s", MSG429, e.err.Error())
	case 500:
		return fmt.Sprintf("%s: %s", MSG500, e.err.Error())
	case 501:
		return fmt.Sprintf("%s: %s", MSG501, e.err.Error())
	}

	return e.err.Error()

}

func NewError(statusCode uint32, msg string, rawErr error, data any) Error {
	return Error{
		status: statusCode,
		msg:    msg,
		err:    rawErr,
		data:   data,
	}
}

func RespError(c *nf.Ctx, err error) error {
	if err == nil {
		return Resp(c, 500, MSG500, "response with nil error", nil)
	}

	var re = &Error{}
	if errors.As(err, re) {
		if re.err == nil {
			return Resp(c, re.status, re.msg, re.msg, re.data)
		}

		return Resp(c, re.status, re.msg, re.err.Error(), re.data)
	}

	return Resp(c, 500, MSG500, err.Error(), nil)
}
