package migration

import (
	"embed"
)

//go:embed migration/*.sql
var EmbedMigrations embed.FS
