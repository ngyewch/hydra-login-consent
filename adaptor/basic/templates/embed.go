package templates

import "embed"

var (
	//go:embed *.gohtml
	TemplateFS embed.FS
)
