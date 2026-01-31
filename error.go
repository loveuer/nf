package ursa

import "strconv"

type Err struct {
	Status int
	Msg    string
}

func (n Err) Error() string {
	return strconv.Itoa(n.Status) + " " + n.Msg
}

func NewNFError(status int, msg string) Err {
	return Err{Status: status, Msg: msg}
}
