package funtemplates

import (
	"bytes"
	"errors"
	"fmt"
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
		{
			name:       "Int to bool test",
			template:   "{{ map $.DataInts $.Funcs.odd }}",
			want:       "[true false true false]",
			correctErr: NoError,
		},
		{
			name:       "No parameters lambda to bool test",
			template:   "{{ map $.DataInts $.Funcs.false }}",
			want:       "[false false false false]",
			correctErr: NoError,
		},
		{
			name:       "Correct error on not a func",
			template:   "{{ map $.DataInts $.InvalidFuncs.NotAFunction }}",
			want:       "",
			correctErr: ErrorIs(ErrExpected2ndArgumentToBeFunction),
		},
		{
			name:       "First parameter must be a slice not a string",
			template:   `{{ map "asdfasdf" $.Funcs.false }}`,
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstParameterToBeSlice),
		},
		{
			name:       "First parameter must be a slice not a number",
			template:   "{{ map 123 $.Funcs.false }}",
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstParameterToBeSlice),
		},
		{
			name:       "First parameter must be a slice. Null is a slice",
			template:   "{{ map nil $.Funcs.false }}",
			want:       "[]",
			correctErr: NoError,
		},
		{
			name:       "First parameter must be a slice not a number",
			template:   "{{ map $.DataInts $.InvalidFuncs.TooManyArgs }}",
			want:       "",
			correctErr: ErrorIs(ErrInputFuncMustTake0or1Arguments),
		},
		{
			name:       "No more than 2 returns",
			template:   "{{ map $.DataInts $.InvalidFuncs.TooManyReturns }}",
			want:       "",
			correctErr: ErrorIs(ErrExpected1Or2ReturnTypes),
		},
		{
			name:       "No more fewer than 1 returns",
			template:   "{{ map $.DataInts $.InvalidFuncs.NoReturns }}",
			want:       "",
			correctErr: ErrorIs(ErrExpected1Or2ReturnTypes),
		},
	}
	funcs := misc.MergeMaps(TextFunctions(), misc.SimpleTextFunctions())
	data := struct {
		DataInts     []int
		Funcs        map[string]any
		InvalidFuncs map[string]any
	}{
		DataInts: []int{1, 2, 3, 4},
		Funcs:    funcs,
		InvalidFuncs: map[string]any{
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpl := template.Must(template.New("").Funcs(funcs).Parse(tt.template))
			got := bytes.NewBuffer(nil)
			err := tmpl.Execute(got, data)
			if tt.correctErr != nil {
				if description, ok := tt.correctErr(err); !ok {
					t.Errorf("MapTemplateFunc() got error =\n> %v\n\n%s", err, description)
					return
				}
				if err != nil {
					return
				}
			}
			if diff := cmp.Diff(tt.want, got.String()); diff != "" {
				t.Errorf("MapTemplateFunc() diff =\n %s", diff)
			}
		})
	}
}

func ErrorIs(shouldBeErr error) func(err error) (string, bool) {
	return func(err error) (string, bool) {
		description := fmt.Sprintf("Expected:\n> %s", "no error")
		if shouldBeErr != nil {
			description = fmt.Sprintf("Expected:\n> %s", shouldBeErr.Error())
		}
		return description, errors.Is(err, shouldBeErr)
	}
}

func NoError(err error) (string, bool) {
	if err == nil {
		return "", true
	}
	return fmt.Sprintf("Expected:\n> %s", "no error"), false
}
