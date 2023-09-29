package public

import (
	"embed"
)

//go:embed web/*
var Assets embed.FS
