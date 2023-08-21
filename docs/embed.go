/*
Package docs embeds the docs folder into the binary.
the docs website will serve at `/docs` path.
*/
package docs

import "embed"

//go:embed *

// FS fs
var FS embed.FS
