package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const everyonePerms = 0o777

//go:embed platforms.go.template
var platformsTemplate []byte

//go:embed definitionoptions.go.template
var definitionOptionsTemplate []byte

var (
	platformPattern = regexp.MustCompile(`const (\w+) =`)
	registerPattern = regexp.MustCompile(`func (register\w+)`)
)

const (
	platformsFileName         = "cli/platforms.go"
	definitionOptionsFileName = "cli/definitionoptions/definitionoptions.go"
)

func convertDefinitionFileNameToGoName(s string) string {
	ss := strings.Split(s, "_")

	out := ""

	for _, part := range ss {
		out = fmt.Sprintf("%s%s%s", out, strings.ToTitle(string(part[0])), part[1:])
	}

	return out
}

func handlePlatforms() {
	t, err := template.New("platforms.go").Parse(string(platformsTemplate))
	if err != nil {
		panic("failed parsing platforms template")
	}

	definitionFiles, err := filepath.Glob("./assets/definitions/*.yaml")
	if err != nil {
		panic("failed globbing definition files")
	}

	definitionFileToGoName := map[string]string{}

	for _, definitionFile := range definitionFiles {
		f := strings.TrimSuffix(filepath.Base(definitionFile), ".yaml")

		definitionFileToGoName[f] = convertDefinitionFileNameToGoName(f)
	}

	var rendered bytes.Buffer

	err = t.Execute(
		&rendered,
		struct {
			PlatformMap map[string]string
		}{
			PlatformMap: definitionFileToGoName,
		},
	)
	if err != nil {
		panic("failed executing template")
	}

	err = os.WriteFile(
		platformsFileName,
		rendered.Bytes(),
		everyonePerms,
	)
	if err != nil {
		panic("failed writing platform template")
	}
}

func handleOptions() {
	t, err := template.New("definitionoptions.go").Parse(string(definitionOptionsTemplate))
	if err != nil {
		panic("failed parsing definitionoptions template")
	}

	definitionOptions, err := filepath.Glob("./cli/definitionoptions/*.go")
	if err != nil {
		panic("failed globbing definition options files")
	}

	definitionOptionsPlatformToRegisterFuncs := map[string]string{}

	for _, defintionOptionsFile := range definitionOptions {
		if defintionOptionsFile == definitionOptionsFileName {
			continue
		}

		definitionsOptionsFileContents, err := os.ReadFile(defintionOptionsFile) //nolint: gosec
		if err != nil {
			panic("failed reading definitions file")
		}

		foundPlatform := platformPattern.FindSubmatch(definitionsOptionsFileContents)
		foundRegisterFunc := registerPattern.FindSubmatch(definitionsOptionsFileContents)

		definitionOptionsPlatformToRegisterFuncs[string(foundPlatform[1])] = string(
			foundRegisterFunc[1],
		)
	}

	var rendered bytes.Buffer

	err = t.Execute(
		&rendered,
		struct {
			OptionMap map[string]string
		}{
			OptionMap: definitionOptionsPlatformToRegisterFuncs,
		},
	)
	if err != nil {
		panic("failed executing template")
	}

	err = os.WriteFile(
		definitionOptionsFileName,
		rendered.Bytes(),
		everyonePerms,
	)
	if err != nil {
		panic("failed writing rendered template")
	}
}

func main() {
	handlePlatforms()
	handleOptions()
}
