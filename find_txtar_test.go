package funtemplates

import (
	"bytes"
	"embed"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/arran4/go-template-functional-operations/misc"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/*.txtar
var testData embed.FS

func TestFindIndexTxtar(t *testing.T) {
	entries, err := testData.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	funcs := misc.MergeMaps(TextFunctions(), misc.SimpleTextFunctions())
	// Add some custom functions for testing if not already present
	funcs["alwaysTrue"] = func() bool { return true }
	funcs["alwaysFalse"] = func() bool { return false }
	funcs["even"] = func(i int) bool { return i%2 == 0 }

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txtar") {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			content, err := testData.ReadFile(filepath.Join("testdata", entry.Name()))
			if err != nil {
				t.Fatal(err)
			}

			archive := txtar.Parse(content)
			var inputData map[string]any
			var tmplStr string
			var want string

			for _, f := range archive.Files {
				switch f.Name {
				case "input.yaml":
					if err := yaml.Unmarshal(f.Data, &inputData); err != nil {
						t.Fatalf("failed to unmarshal input.yaml: %v", err)
					}
				case "template.tmpl":
					tmplStr = string(f.Data)
				case "output.txt":
					want = strings.TrimSpace(string(f.Data))
				}
			}

			// Add functions to the input data so they can be accessed via .Funcs
			if inputData == nil {
				inputData = make(map[string]any)
			}
			inputData["Funcs"] = funcs

			tmpl, err := template.New(entry.Name()).Funcs(funcs).Parse(tmplStr)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var got bytes.Buffer
			if err := tmpl.Execute(&got, inputData); err != nil {
				t.Fatalf("failed to execute template: %v", err)
			}

			if diff := cmp.Diff(want, strings.TrimSpace(got.String())); diff != "" {
				t.Errorf("template execution mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
