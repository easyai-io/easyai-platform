/*
Package www embeds the www directory into the binary.
the www website will serve at `/` path.
*/
package www

import "embed"

//go:embed *

// FS fs
var FS embed.FS
