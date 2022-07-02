package assets

import "embed"

//go:embed platforms/*
// Assets is the embedded assets objects for the included platform yaml data.
var Assets embed.FS
