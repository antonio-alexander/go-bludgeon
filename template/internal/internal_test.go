package internal_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/antonio-alexander/go-bludgeon/template/internal"

	"github.com/stretchr/testify/assert"
)

func parse(packageName string, objects []internal.Object) (string, error) {
	var s []string

	p, err := internal.Execute(internal.TemplatePackage, internal.Package{
		Name: packageName,
	})
	if err != nil {
		return "", err
	}
	s = append(s, p)
	for _, object := range objects {
		o, err := internal.Execute(internal.TemplateObject, object)
		if err != nil {
			return "", err
		}
		s = append(s, o)
	}
	return strings.TrimSpace(strings.Join(s, "\n")), nil
}

func TestObject(t *testing.T) {
	objects := []internal.Object{
		{
			Name:    "Sample",
			Comment: "Sample is a big bad sample",
			Fields: []internal.Field{
				{
					Name:      "Hello",
					Tag:       "hello",
					Type:      "string",
					Comment:   "does absolutely nothing",
					JSON:      true,
					OmitEmpty: false,
				},
				{
					Name:      "World",
					Tag:       "world",
					Type:      "bool",
					Comment:   "does absolutely nothing second",
					JSON:      true,
					OmitEmpty: false,
				},
			},
		},
		{
			Name:    "SamplePartial",
			Comment: "SamplePartial is a big partial sample",
			Fields: []internal.Field{
				{
					Name:    "Hello",
					Tag:     "hello",
					Type:    "*string",
					Comment: "Hello does absolutely nothing",
					JSON:    true,
					YAML:    true,
				},
			},
		},
	}
	s, err := parse("data", objects)
	assert.Nil(t, err)
	assert.NotEmpty(t, s)
	fmt.Println(s)
}
