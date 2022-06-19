package internal

import (
	"bytes"
	"text/template"
)

func Execute(t *template.Template, item interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, item); err != nil {
		return "", err
	}
	return buf.String(), nil
}
