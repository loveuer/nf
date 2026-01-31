package resp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/loveuer/ursa"
)

func handleEmptyMsg(status uint32, msg string) string {
	if msg == "" {
		switch status {
		case 200:
			msg = MSG200
		case 202:
			msg = MSG202
		case 400:
			msg = MSG400
		case 401:
			msg = MSG401
		case 403:
			msg = MSG403
		case 404:
			msg = MSG404
		case 429:
			msg = MSG429
		case 500:
			msg = MSG500
		case 501:
			msg = MSG501
		}
	}

	return msg
}

func Resp(c *ursa.Ctx, status uint32, msg string, err string, data any) error {
	msg = handleEmptyMsg(status, msg)

	c.Set(RealStatusHeader, strconv.Itoa(int(status)))

	if data == nil {
		return c.JSON(ursa.Map{"status": status, "msg": msg, "err": err})
	}

	return c.JSON(ursa.Map{"status": status, "msg": msg, "err": err, "data": data})
}

func Resp200(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG200

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 200, msg, "", data)
}

func Resp202(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG202

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 202, msg, "", data)
}

func Resp400(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG400
	err := ""

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
		err = msg
	}

	return Resp(c, 400, msg, err, data)
}

func Resp401(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG401
	err := ""

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
		err = msg
	}

	return Resp(c, 401, msg, err, data)
}

func Resp403(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG403
	err := ""

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
		err = msg
	}

	return Resp(c, 403, msg, err, data)
}

func Resp418(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG418
	err := ""

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
		err = ""
	}

	return Resp(c, 418, msg, err, data)
}

func Resp429(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG429
	err := ""

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
		err = ""
	}

	return Resp(c, 429, msg, err, data)
}

func Resp500(c *ursa.Ctx, data any, msgs ...string) error {
	msg := MSG500
	err := ""

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
		err = msg
	}

	return Resp(c, 500, msg, err, data)
}
