package tp

import (
	"bytes"
	"text/template"
)

var (
	_t *template.Template
)

func init() {
	_t = template.New("tp")
}

func RenderVar(t []byte, data map[string]any) ([]byte, error) {
	tr, err := _t.Parse(string(t))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err = tr.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
