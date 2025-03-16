//go:build linux && amd64 && musl
// +build linux,amd64,musl

package assets

import "embed"

// Lib is the embedded libscrapli shared object.
//
//go:embed lib/x86_64-linux-musl
var Lib embed.FS
