package funtemplates

import (
	"bytes"
	"embed"
	"encoding/json"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/arran4/go-template-functional-operations/misc"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

//go:embed testdata/filter.txt
var filterTestData embed.FS

type TxtarTestInput struct {
	DataInts []int
	DataAny  []any
}

func TestFilterTxtar(t *testing.T) {
	archiveData, err := filterTestData.ReadFile("testdata/filter.txt")
	if err != nil {
		t.Fatal(err)
	}
	archive := txtar.Parse(archiveData)

	// Group files by directory (test case)
	testCases := make(map[string]map[string]string)
	var orderedNames []string

	for _, f := range archive.Files {
		dir, base := path.Split(f.Name)
		dir = strings.TrimRight(dir, "/")
		if dir == "" {
			continue
		}
		if _, ok := testCases[dir]; !ok {
			testCases[dir] = make(map[string]string)
			orderedNames = append(orderedNames, dir)
		}
		testCases[dir][base] = string(f.Data)
	}

	funcs := misc.MergeMaps(TextFunctions(), misc.SimpleTextFunctions())

	for _, name := range orderedNames {
		files := testCases[name]
		t.Run(name, func(t *testing.T) {
			inputJson := files["input.json"]
			tmplStr := files["template.tmpl"]
			want := strings.TrimSpace(files["output.txt"])

			var data TxtarTestInput
			if inputJson != "" {
				if err := json.Unmarshal([]byte(inputJson), &data); err != nil {
					t.Fatalf("failed to unmarshal input: %v", err)
				}
			}

			// Fixup DataAny to contain ints instead of float64s (from JSON)
			// This allows us to test []interface{} containing ints, which matches the behavior of code calling the library.
			for i, v := range data.DataAny {
				if f, ok := v.(float64); ok {
					data.DataAny[i] = int(f)
				}
			}

			// Prepare data map including Funcs
			tmplData := struct {
				DataInts []int
				DataAny  []any
				Funcs    map[string]any
			}{
				DataInts: data.DataInts,
				DataAny:  data.DataAny,
				Funcs:    funcs,
			}

			tmpl, err := template.New("").Funcs(funcs).Parse(tmplStr)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, tmplData); err != nil {
				t.Fatalf("failed to execute template: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
