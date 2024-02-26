package web

import "embed"

//go:generate tailwindcss build -i ./assets/base.css -o assets/tailwind.css

//go:embed assets
var Assets embed.FS

//go:embed templates
var TemplatesFS embed.FS
