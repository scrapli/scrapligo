package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/netconf"
)

func main() {

	d, _ := netconf.NewNetconfDriver(
		"localhost",
		base.WithPort(21830),
		base.WithAuthStrictKey(false),
		base.WithAuthUsername("vrnetlab"),
		base.WithAuthPassword("VR-netlab9"),
	)

	err := d.Open()
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)
		return
	}
	defer d.Close()

	r, err := d.GetConfig("running")
	if err != nil {
		fmt.Printf("failed to get config; error: %+v\n", err)
		return
	}

	fmt.Printf("Get Config Response:\n%s\n", r.Result)

	filter := "" +
		"<interfaces xmlns=\"urn:ietf:params:xml:ns:yang:ietf-interfaces\">\n" +
		"  <interface>\n" +
		"    <name>\n" +
		"      GigabitEthernet1\n" +
		"    </name>\n" +
		"  </interface>\n" +
		"</interfaces>"
	r, err = d.Get(netconf.WithNetconfFilter(filter))
	if err != nil {
		fmt.Printf("failed to get with filter; error: %+v\n", err)
		return
	}

	fmt.Printf("Get Response: %s\n", r.Result)

	edit := "" +
		"<config>\n" +
		"    <interfaces xmlns=\"urn:ietf:params:xml:ns:yang:ietf-interfaces\">\n" +
		"        <interface>\n" +
		"            <name>GigabitEthernet1</name>\n" +
		"            <description>scrapliGO was here!</description>\n" +
		"        </interface>\n" +
		"    </interfaces>\n" +
		"</config>"
	r, err = d.EditConfig("running", edit)
	if err != nil {
		fmt.Printf("failed to edit config; error: %+v\n", err)
		return
	}

	fmt.Printf("Edit Config Response: %s\n", r.Result)

	_, _ = d.Commit()
}
