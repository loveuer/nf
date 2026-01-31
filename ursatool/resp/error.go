package resp

import (
	"errors"
	"github.com/loveuer/ursa"
)

type Error struct {
	status uint32
	msg    string
	err    error
	data   any
}

func (e Error) Error() string {
	if e.msg != "" {
		return e.msg
	}

	switch e.status {
	case 200:
		return MSG200
	case 202:
		return MSG202
	case 400:
		return MSG400
	case 401:
		return MSG401
	case 403:
		return MSG403
	case 404:
		return MSG404
	case 429:
		return MSG429
	case 500:
		return MSG500
	case 501:
		return MSG501
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

func RespError(c *ursa.Ctx, err error) error {
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
