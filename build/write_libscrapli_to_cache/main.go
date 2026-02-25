package main

import (
	"context"
	"fmt"

	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
)

// just a simple program to easily expose the EnsureLibscrapli function for devs. obv just
// `go run build/write_libscrapli_to_cache.go`.
func main() {
	p, err := scrapligoffi.EnsureLibscrapli(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("libscrapli is available at path %q\n", p) //nolint: forbidigo
}
