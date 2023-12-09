package static

import "embed"

var (
	//go:embed images
	StaticFS embed.FS
)
