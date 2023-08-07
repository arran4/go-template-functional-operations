package funtemplates

import (
	ht "text/template"
	tt "text/template"
)

func TextFunctions() tt.FuncMap {
	return map[string]any{
		"filter":    FilterTemplateFunc,
		"find":      FindTemplateFunc,
		"findIndex": FindIndexTemplateFunc,
		"map":       MapTemplateFunc,
	}
}

func HtmlFunctions() ht.FuncMap {
	return map[string]any{
		"filter":    FilterTemplateFunc,
		"find":      FindTemplateFunc,
		"findIndex": FindIndexTemplateFunc,
		"map":       MapTemplateFunc,
	}
}
