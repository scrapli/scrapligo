module github.com/scrapli/scrapligo

go 1.16

require (
	github.com/creack/pty v1.1.21
	github.com/google/go-cmp v0.6.0
	github.com/sirikothe/gotextfsm v1.0.1-0.20200816110946-6aa2cfd355e4
	// v0.7.0 is the latest we can be on if we want to support 1.16
	// and v0.6.0 if we want to not piss off linters it seems
	golang.org/x/crypto v0.6.0
	gopkg.in/yaml.v3 v3.0.1
)
