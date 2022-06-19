package internal

import (
	"text/template"
)

const (
	templateVariable  string = templateComment + "var {{ .Name }} = {{ .Value }}"
	templateVariables string = templateComment + " {{ .Name }} = {{ .Value }}"
)

type Variable struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Value   string `json:"value"`
}

var (
	TemplateVariable  = template.Must(template.New("variable").Parse(templateVariable))
	TemplateVariables = template.Must(template.New("variables").Parse(templateVariables))
)
