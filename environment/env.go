package environment

import "html/template"

type Env struct {
	Templates map[string]*template.Template
}
