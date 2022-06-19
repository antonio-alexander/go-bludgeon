package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"
)

type Todo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

const templateStruct string = "//{{ .Description}}\ntype {{ .Name}} struct {}\n"

func main() {
	var sample Todo

	byts, err := ioutil.ReadFile("./sample.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(byts, &sample); err != nil {
		panic(err)
	}
	t, err := template.New("todos").Parse(templateStruct)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, sample)
	if err != nil {
		panic(err)
	}
	fmt.Print(buf.String())
}
