package main

import (
	"fmt"
	"github.com/arran4/go-template-functional-operations"
	"github.com/arran4/go-template-functional-operations/misc"
	"os"
	"text/template"
)

func main() {
	funcs := map[string]any{
		"inc": func(i int) int {
			return i + 1
		},
		"odd": func(i int) bool {
			return i%2 == 1
		},
	}
	funcs = misc.MergeMaps(funtemplates.TextFunctions(), funcs)

	tmpl := template.Must(template.New("").Funcs(funcs).Parse(`
        {{ map $.Data $.Funcs.inc }}
        {{ map $.Data $.Funcs.odd }}
    `))

	data := struct {
		Data  []int
		Funcs map[string]any
	}{
		Data:  []int{1, 2, 3, 4},
		Funcs: funcs,
	}

	err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println(err)
	}
}
