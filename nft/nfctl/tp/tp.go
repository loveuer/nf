package tp

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type TpCmd interface {
	String() string
	Execute() error
}

var (
	_ TpCmd = (*TpGenerate)(nil)
	_ TpCmd = (*TpReplace)(nil)
)

type TpGenerate struct {
	pwd      string
	filename string
	content  []string
}

func (t *TpGenerate) String() string {
	//TODO implement me
	panic("implement me")
}

func (t *TpGenerate) Execute() error {
	var (
		err      error
		location = t.filename
		input    *os.File
	)

	if !path.IsAbs(t.filename) {
		location = path.Join(t.pwd, t.filename)
	}

	if err = os.MkdirAll(path.Dir(location), 0644); err != nil {
		return err
	}

	if !strings.HasSuffix(location, "/") {
		if input, err = os.OpenFile(location, os.O_CREATE|os.O_APPEND, 0744); err != nil {
			return err
		}

		if len(t.content) > 0 {
			content := strings.Join(t.content, "\n")
			_, err = input.WriteString(content)
			return err
		}
	}

	return nil
}

func newGenerate(pwd string, lines []string) (*TpGenerate, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("generate cmd require file/folder name")
	}

	return &TpGenerate{
		pwd:      pwd,
		filename: lines[0],
		content:  lines[1:],
	}, nil
}

type TpReplace struct {
	pwd string
}

func (t *TpReplace) String() string {
	//TODO implement me
	panic("implement me")
}

func (t *TpReplace) Execute() error {
	//TODO implement me
	panic("implement me")
}

func newReplace(pwd string, lines []string) (*TpReplace, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid replace cmd")
	}
	return &TpReplace{pwd: pwd}, nil
}
