//go:build darwin && arm64
// +build darwin,arm64

package assets

import "embed"

// Lib is the embedded libscrapli shared object.
//
//go:embed lib/aarch64-macos
var Lib embed.FS
