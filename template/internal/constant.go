package internal

import (
	"text/template"
)

const (
	templateConstant  string = templateComment + "const {{ .Name }} = {{ .Value }}"
	templateConstants string = templateComment + " {{ .Name }} = {{ .Value }}"
)

type Constant struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Value   string `json:"value"`
}

var (
	TemplateConstant  = template.Must(template.New("constant").Parse(templateConstant))
	TemplateConstants = template.Must(template.New("constants").Parse(templateConstants))
)
