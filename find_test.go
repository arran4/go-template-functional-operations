package funtemplates

import (
	"bytes"
	"github.com/arran4/go-template-functional-operations/misc"
	"github.com/google/go-cmp/cmp"
	"testing"
	"text/template"
)

func TestFindTemplateFunc(t *testing.T) {
	tests := []struct {
		name       string
		template   string
		want       string
		correctErr func(err error) (string, bool)
	}{
		{
			name:       "Int returns error",
			template:   "{{ find $.DataInts $.Funcs.inc }}",
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstReturnToBeBool),
		},
		{
			name:       "first Odd",
			template:   "{{ find $.DataInts $.Funcs.odd }}",
			want:       "1",
			correctErr: NoError,
		},
		{
			name:       "False causes find to return nil",
			template:   "{{ find $.DataInts $.Funcs.false }}",
			want:       "<no value>",
			correctErr: NoError,
		},
		{
			name:       "Correct error on not a func",
			template:   "{{ find $.DataInts $.InvalidFuncs.NotAFunction }}",
			want:       "",
			correctErr: ErrorIs(ErrExpected2ndArgumentToBeFunction),
		},
		{
			name:       "First parameter must be a slice not a string",
			template:   `{{ find "asdfasdf" $.Funcs.false }}`,
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstParameterToBeSlice),
		},
		{
			name:       "First parameter must be a slice not a number",
			template:   "{{ find 123 $.Funcs.false }}",
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstParameterToBeSlice),
		},
		{
			name:       "First parameter must be a slice. Null is a slice",
			template:   "{{ find nil $.Funcs.false }}",
			want:       "<no value>",
			correctErr: NoError,
		},
		{
			name:       "First parameter must be a slice not a number",
			template:   "{{ find $.DataInts $.InvalidFuncs.TooManyArgs }}",
			want:       "",
			correctErr: ErrorIs(ErrInputFuncMustTake0or1Arguments),
		},
		{
			name:       "No more than 2 returns",
			template:   "{{ find $.DataInts $.InvalidFuncs.TooManyReturns }}",
			want:       "",
			correctErr: ErrorIs(ErrExpected1Or2ReturnTypes),
		},
		{
			name:       "No more fewer than 1 returns",
			template:   "{{ find $.DataInts $.InvalidFuncs.NoReturns }}",
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
			"NoReturns": func() {
				return
			},
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
func TestFindIndexTemplateFunc(t *testing.T) {
	tests := []struct {
		name       string
		template   string
		want       string
		correctErr func(err error) (string, bool)
	}{
		{
			name:       "Int returns error",
			template:   "{{ findIndex $.DataInts $.Funcs.inc }}",
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstReturnToBeBool),
		},
		{
			name:       "first Odd",
			template:   "{{ findIndex $.DataInts $.Funcs.odd }}",
			want:       "0",
			correctErr: NoError,
		},
		{
			name:       "False causes findIndex to return -1",
			template:   "{{ findIndex $.DataInts $.Funcs.false }}",
			want:       "-1",
			correctErr: NoError,
		},
		{
			name:       "Correct error on not a func",
			template:   "{{ findIndex $.DataInts $.InvalidFuncs.NotAFunction }}",
			want:       "",
			correctErr: ErrorIs(ErrExpected2ndArgumentToBeFunction),
		},
		{
			name:       "First parameter must be a slice not a string",
			template:   `{{ findIndex "asdfasdf" $.Funcs.false }}`,
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstParameterToBeSlice),
		},
		{
			name:       "First parameter must be a slice not a number",
			template:   "{{ findIndex 123 $.Funcs.false }}",
			want:       "",
			correctErr: ErrorIs(ErrExpectedFirstParameterToBeSlice),
		},
		{
			name:       "First parameter must be a slice. Null is a slice",
			template:   "{{ findIndex nil $.Funcs.false }}",
			want:       "-1",
			correctErr: NoError,
		},
		{
			name:       "First parameter must be a slice not a number",
			template:   "{{ findIndex $.DataInts $.InvalidFuncs.TooManyArgs }}",
			want:       "",
			correctErr: ErrorIs(ErrInputFuncMustTake0or1Arguments),
		},
		{
			name:       "No more than 2 returns",
			template:   "{{ findIndex $.DataInts $.InvalidFuncs.TooManyReturns }}",
			want:       "",
			correctErr: ErrorIs(ErrExpected1Or2ReturnTypes),
		},
		{
			name:       "No more fewer than 1 returns",
			template:   "{{ findIndex $.DataInts $.InvalidFuncs.NoReturns }}",
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
			"NoReturns": func() {
				return
			},
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
