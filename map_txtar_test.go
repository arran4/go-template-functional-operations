package funtemplates

import (
	"bytes"
	"embed"
	"encoding/json"
	"github.com/arran4/go-template-functional-operations/misc"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
	"io/fs"
	"testing"
	"text/template"
)

//go:embed testdata/map/*.txtar
var mapTestData embed.FS

func TestMapTemplateFunc_Txtar(t *testing.T) {
	files, err := fs.ReadDir(mapTestData, "testdata/map")
	if err != nil {
		t.Fatal(err)
	}

	funcs := misc.MergeMaps(TextFunctions(), misc.SimpleTextFunctions())
	invalidFuncs := map[string]any{
		"NotAFunction": "this is totally not a function",
		"TooManyArgs": func(one, two, three, four, five int) int {
			return 0
		},
		"TooManyReturns": func() (int, int, int) {
			return 0, 1, 3
		},
		"NotAnError": func() (int, int) {
			return 0, 2
		},
		"NoReturns": func() {},
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			content, err := mapTestData.ReadFile("testdata/map/" + file.Name())
			if err != nil {
				t.Fatal(err)
			}
			archive := txtar.Parse(content)

			var inputJSON []byte
			var tmplBody string
			var expect string

			for _, f := range archive.Files {
				switch f.Name {
				case "input.json":
					inputJSON = f.Data
				case "template.tmpl":
					tmplBody = string(f.Data)
				case "expect.txt":
					expect = string(f.Data)
				}
			}

			// Define the data structure matching existing tests
			data := struct {
				DataInts     []int
				Funcs        map[string]any
				InvalidFuncs map[string]any
			}{
				Funcs:        funcs,
				InvalidFuncs: invalidFuncs,
			}

			// Unmarshal input if present to override defaults or add data
			if len(inputJSON) > 0 {
				if err := json.Unmarshal(inputJSON, &data); err != nil {
					t.Fatalf("failed to unmarshal input.json: %v", err)
				}
			}

			tmpl, err := template.New("").Funcs(funcs).Parse(tmplBody)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			out := bytes.NewBuffer(nil)
			err = tmpl.Execute(out, data)

			got := out.String()
			// If execution failed, check if we expected an error
			if err != nil {
				got = err.Error()
			}

            // Trim newlines for comparison convenience if needed,
            // but txtar usually preserves them. Let's be exact first.
            // If expected output has a trailing newline (from txtar file), we might need to handle it.
            // txtar.Parse leaves content as bytes.

            // Simple trim for safety if users edit files with newlines
            if diff := cmp.Diff(trim(expect), trim(got)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func trim(s string) string {
	return string(bytes.TrimSpace([]byte(s)))
}
