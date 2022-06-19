package internal

import (
	"text/template"
)

const (
	templateTagf         string = " `" + templateTagJSON + templateTagYAML + "`"
	templateTagJSON      string = "{{if .JSON}}json:\"{{ .Tag}}" + templateTagOmitEmpty + "\"{{end}}"
	templateTagYAML      string = "{{if .YAML}} yaml:\"{{ .Tag}}" + templateTagOmitEmpty + "\"{{end}}"
	templateTagOmitEmpty string = "{{if .OmitEmpty}},omitempty{{end}}"
	templateField        string = `{{range .Fields}}` + templateComment + `{{printf "\n%s %s" .Name .Type}}` + templateTagf + `{{end}}`
	templateStruct       string = templateComment + `{{printf "\ntype %s struct {" .Name }}` + templateField + "\n}"
)

type Field struct {
	Name      string `json:"name"`
	Tag       string `json:"name_tag"`
	Type      string `json:"type"`
	Comment   string `json:"comment"`
	JSON      bool   `json:"json"`
	YAML      bool   `json:"yaml"`
	OmitEmpty bool   `json:"omitempty"`
}

type Object struct {
	Name    string  `json:"name"`
	Comment string  `json:"comment"`
	Fields  []Field `json:"fields"`
}

var TemplateObject = template.Must(template.New("object").Parse(templateStruct))
