package berrypost

import (
	"embed"
)

//go:embed statics/templates/*
var TemplateFS embed.FS

//go:embed statics/dist/*
var DistFS embed.FS
