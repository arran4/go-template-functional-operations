package main

import (
	"fmt"
	"github.com/arran4/go-template-functional-operations"
	"os"
	"text/template"
)

func main() {
	// Create a template.
	funcs := map[string]any{
		"map": funtemplates.MapTemplateFunc,
		"inc": func(i int) int {
			return i + 1
		},
		"odd": func(i int) bool {
			return i%2 == 1
		},
	}
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(`
        {{ map $.Data $.Funcs.inc }}
        {{ map $.Data $.Funcs.odd }}
    `))

	// Create some data to be used in the template.
	data := struct {
		Data  []int
		Funcs map[string]any
	}{
		Data:  []int{1, 2, 3, 4},
		Funcs: funcs,
	}

	// Execute the template with the data.
	err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println(err)
	}
}
