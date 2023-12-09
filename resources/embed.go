package resources

import "embed"

var (
	//go:embed templates
	TemplateFS embed.FS
)
