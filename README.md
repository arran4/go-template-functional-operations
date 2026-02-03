# Go Template Functional Operations

[![Go Report Card](https://goreportcard.com/badge/github.com/arran4/go-template-functional-operations)](https://goreportcard.com/report/github.com/arran4/go-template-functional-operations)
[![GoDoc](https://godoc.org/github.com/arran4/go-template-functional-operations?status.svg)](https://godoc.org/github.com/arran4/go-template-functional-operations)

This library adds reflection-based functional programming primitives to Go's standard `text/template` and `html/template` packages.

It currently supports:
*   `map`
*   `filter`
*   `find`
*   `findIndex`

## Why use this?

Go's standard template library is intentionally minimal to encourage separation of logic. However, strict separation can sometimes lead to boilerplate code where you create specific View Models just to handle minor display logic (e.g., formatting a list, filtering out empty values, or finding a specific item).

This library bridges that gap by allowing you to perform common data transformations and queries directly within your templates, keeping your Go code cleaner and your templates more powerful.

## Installation

```bash
go get github.com/arran4/go-template-functional-operations
```

## Usage

### 1. Register the functions

You need to register the functions with your template before parsing it. The library provides `TextFunctions()` and `HtmlFunctions()` for this purpose.

### 2. Prepare your data

**Important:** The functional operations (`map`, `filter`, etc.) expect a **function** as their second argument. In Go templates, you cannot easily pass a function registered in `FuncMap` as an argument to another function.

Therefore, you must pass your helper/predicate functions as part of your data structure (or a map) that you execute the template with.

### Example

```go
package main

import (
	"os"
	"text/template"

	funtemplates "github.com/arran4/go-template-functional-operations"
)

func main() {
    // Define the functions you want to use inside the template
    myFuncs := map[string]any{
        "inc": func(i int) int { return i + 1 },
        "isOdd": func(i int) bool { return i%2 != 0 },
        "isThree": func(i int) bool { return i == 3 },
    }

    // Prepare your data.
    // We pass our functions in the 'F' field so we can access them in the template.
    data := struct {
        Items []int
        F     map[string]any
    }{
        Items: []int{1, 2, 3, 4, 5},
        F:     myFuncs,
    }

    // Create template and register the library's functions
    t := template.New("demo").Funcs(funtemplates.TextFunctions())

    // Parse template
    src := `
    Original:  {{ .Items }}
    Mapped:    {{ map .Items .F.inc }}
    Filtered:  {{ filter .Items .F.isOdd }}
    Find:      {{ find .Items .F.isThree }}
    FindIndex: {{ findIndex .Items .F.isThree }}
    `

    template.Must(t.Parse(src))

    if err := t.Execute(os.Stdout, data); err != nil {
        panic(err)
    }
}
```

**Output:**
```
    Original:  [1 2 3 4 5]
    Mapped:    [2 3 4 5 6]
    Filtered:  [1 3 5]
    Find:      3
    FindIndex: 2
```

## API Reference

### `map`

Applies a function to every element in a slice and returns a new slice with the results.

*   **Signature:** `func(slice any, f any) (any, error)`
*   **Arguments:**
    *   `slice`: The input slice (any type).
    *   `f`: A function that takes one argument (element of slice) and returns a value (and optional error).
*   **Returns:** A new slice containing the results of applying `f` to each element.

### `filter`

Iterates over a slice and returns a new slice containing only the elements for which the predicate function returns `true`.

*   **Signature:** `func(slice any, f any) (any, error)`
*   **Arguments:**
    *   `slice`: The input slice (any type).
    *   `f`: A predicate function that takes one argument (element of slice) and returns a `bool` (and optional error).
*   **Returns:** A new slice with elements that satisfied the predicate.

### `find`

Returns the first element in the slice that satisfies the provided predicate function.

*   **Signature:** `func(slice any, f any) (any, error)`
*   **Arguments:**
    *   `slice`: The input slice.
    *   `f`: A predicate function returning `bool`.
*   **Returns:** The first matching element, or `nil` if no match is found.

### `findIndex`

Returns the index of the first element in the slice that satisfies the provided predicate function.

*   **Signature:** `func(slice any, f any) (int, error)`
*   **Arguments:**
    *   `slice`: The input slice.
    *   `f`: A predicate function returning `bool`.
*   **Returns:** The index of the first match, or `-1` if no match is found.

## Error Handling

The functions will return an error if:
*   The first argument is not a slice.
*   The second argument is not a function.
*   The function argument does not match the expected signature (e.g., wrong number of arguments or return values).
*   The function itself returns an error (if it supports returning an error).
