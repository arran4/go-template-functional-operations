package funtemplates

import (
	ht "text/template"
	tt "text/template"
)

func TextFunctions() tt.FuncMap {
	return map[string]any{
		"filter": FilterTemplateFunc,
		"map":    MapTemplateFunc,
	}
}

func HtmlFunctions() ht.FuncMap {
	return map[string]any{
		"filter": FilterTemplateFunc,
		"map":    MapTemplateFunc,
	}
}
