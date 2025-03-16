//go:build darwin && amd64
// +build darwin,amd64

package assets

import "embed"

// Lib is the embedded libscrapli shared object.
//
//go:embed lib/x86_64-macos
var Lib embed.FS
