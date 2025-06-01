package main

import (
	"fmt"

	scrapligoffi "github.com/scrapli/scrapligo/ffi"
)

// just a simple program to easily expose the EnsureLibscrapli function for devs. obv just
// `go run build/write_libscrapli_to_cache.go`.
func main() {
	p, err := scrapligoffi.EnsureLibscrapli()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("libscrapli is available at path %q", p)) //nolint: forbidigo
}
