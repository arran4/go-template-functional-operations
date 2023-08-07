# Go Template Function Operations

The goal of this library is to quickly add some reflection based functional functions to the `text/template` and 
`html/template` language. Namely:
* `map`
* `filter`
* `find`
* `findIndex`

This library exists in lieu of generic support in `text/template` or `html/template`.

# Exported functions:

* `func MapTemplateFunc(slice any, f any) (any, error)` The map function.
* `func FilterTemplateFunc(slice any, f any) (any, error)` The map function.
* `func TextFunctions() text/template.FuncMap`
* `func HtmlFunctions() html/template.FuncMap`

# Template Function definitions

## `map`

In go: `MapTemplateFunc`, provided as `map` by `TextFunctions` and `HtmlFunctions`

Definition:
```
func MapTemplateFunc(slice any, f any) (any, error)
```

* The first argument `slice` must be a slice (or nil), of any type.
* Second argument `f` must be a function of these definitions:
  * `func () any` 
  * `func () (any, error)` 
  * `func (v any) any` 
  * `func (v any) (any, error)` 

The return will be:
* The first result: an array of the same length as `slice` in the case, or nil if there was an error. The output of `f` will be in the appropriate place for each value
* The 2nd result: an error if there was an error: See [errors.go](errors.go) for a complete list.

Usage:
* `{{ map $.Data $.Funcs.inc }}`
* `{{ map $.Data $.Funcs.odd }}`

## `filter`

In go: `FilterTemplateFunc`, provided as `map` by `TextFunctions` and `HtmlFunctions`

Definition:
```
func FilterTemplateFunc(slice any, f any) (any, error)
```

* The first argument `slice` must be a slice (or nil), of any type.
* Second argument `f` must be a function of these definitions:
  * `func () bool` 
  * `func () (bool, error)` 
  * `func (v any) bool` 
  * `func (v any) (bool, error)` 

The return will be:
* The first result: an array of the same length or smaller as `slice` in the case, or nil if there was an error.
The output of `f` will be a slice which only contains the values which `f` returned true for `slice`
* The 2nd result: an error if there was an error: See [errors.go](errors.go) for a complete list.

Usage:
* `{{ filter $.Data $.Funcs.odd }}`


# Usage:

```go
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
```

Returns:
```

        [2 3 4 5]
        [true false true false]
    
```
