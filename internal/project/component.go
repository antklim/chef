package project

import "text/template"

type component struct {
	loc      string
	name     string
	template *template.Template
}

func makeComponents(category, server string) map[string]component {
	return nil
}
