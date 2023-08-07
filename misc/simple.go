package misc

import (
	ht "text/template"
	tt "text/template"
)

func SimpleTextFunctions() tt.FuncMap {
	return map[string]any{
		"inc":   incTemplateFunc,
		"odd":   oddTemplateFunc,
		"false": falseTemplateFunc,
	}
}

func SimpleHtmlFunctions() ht.FuncMap {
	return map[string]any{
		"inc": incTemplateFunc,
		"odd": oddTemplateFunc,
	}
}

func incTemplateFunc(i int) int {
	return i + 1
}

func oddTemplateFunc(i int) bool {
	return i%2 == 1
}

func falseTemplateFunc(i int) bool {
	return false
}
