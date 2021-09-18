package main

import (
	"fmt"
	"io/ioutil"

	"github.com/scrapli/scrapligo/cfg"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
)

//Pushes config from a file to a file on device and
//loads it to the device. Cleans up afterwords.
func main() {
	d, err := core.NewCoreDriver(
		"hostname",
		"juniper_junos",
		base.WithPort(22),
		base.WithAuthUsername("root"),
		base.WithAuthPassword("password"),
	)

	if err != nil {
		fmt.Printf("failed to create driver; error: %+v\n", err)
		return
	}

	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)
		return
	}
	defer d.Close()

	c, err := cfg.NewCfgDriver(
		d,
		"juniper_junos",
	)
	if err != nil {
		fmt.Printf("failed to create cfg driver; error: %+v\n", err)
		return
	}

	b, err := ioutil.ReadFile("sample.conf")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(string(b))

	prepareErr := c.Prepare()
	if prepareErr != nil {
		fmt.Printf("failed running prepare method: %v", prepareErr)
	}

	_, err = c.LoadConfig(
		string(b),
		false, //don't load replace. Load merge/set instead
	)
	if err != nil {
		fmt.Printf("failed to load config; error: %+v\n", err)
		return
	}

	_, err = c.CommitConfig()
	if err != nil {
		fmt.Printf("failed to commit config; error: %+v\n", err)
		return
	}
	fmt.Printf("Done loading config\n")
}
