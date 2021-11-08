package verify

import "embed"

// FS is the static asset filesystem.
//go:embed dist
var FS embed.FS
