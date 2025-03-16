//go:build linux && arm64
// +build linux,arm64

package assets

import "embed"

// Lib is the embedded libscrapli shared object.
//
//go:embed lib/aarch64-linux
var Lib embed.FS
