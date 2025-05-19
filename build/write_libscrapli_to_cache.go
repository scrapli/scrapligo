package main

import (
	"fmt"

	scrapligoffi "github.com/scrapli/scrapligo/ffi"
)

func main() {
	p, err := scrapligoffi.EnsureLibscrapli()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("libscrapli is available at path %q", p)) //nolint: forbidigo
}
