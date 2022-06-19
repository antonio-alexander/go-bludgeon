package internal

import (
	"text/template"
)

const templatePackage string = "package {{ .Name}}"

type Package struct {
	Name string `json:"name"`
}

var TemplatePackage = template.Must(template.New("package").Parse(templatePackage))
