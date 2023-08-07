package funtemplates

import (
	"bytes"
	"github.com/arran4/go-template-functional-operations/misc"
	"github.com/google/go-cmp/cmp"
	"testing"
	"text/template"
)

func TestMapTemplateFunc(t *testing.T) {
	tests := []struct {
		name       string
		template   string
		want       string
		correctErr func(err error) (string, bool)
	}{
		{
			name:       "Int to Int test",
			template:   "{{ map $.DataInts $.Funcs.inc }}",
			want:       "[2 3 4 5]",
			correctErr: NoError,
		},
	}
	funcs := misc.MergeMaps(TextFunctions(), misc.SimpleTextFunctions())
	data := struct {
		DataInts []int
		Funcs    map[string]any
	}{
		DataInts: []int{1, 2, 3, 4},
		Funcs:    funcs,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpl := template.Must(template.New("").Funcs(funcs).Parse(tt.template))
			got := bytes.NewBuffer(nil)
			err := tmpl.Execute(got, data)
			if tt.correctErr != nil {
				if description, ok := tt.correctErr(err); !ok {
					t.Errorf("MapTemplateFunc() got error = %v, expected: %v", err, description)
					return
				}
			}
			if diff := cmp.Diff(tt.want, got.String()); diff != "" {
				t.Errorf("MapTemplateFunc() diff =\n %s", diff)
			}
		})
	}
}

func NoError(err error) (string, bool) {
	if err == nil {
		return "", true
	}
	return "No error", false
}
