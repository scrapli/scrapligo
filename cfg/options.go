package cfg

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/scrapli/scrapligo/driver/network"
)

var ErrIgnoredOption = errors.New("option ignored, for different instance type")
var ErrInvalidPlatformAttr = errors.New("invalid platform attribute")

// Option function to set cfg platform options.
type Option func(interface{}) error

// base Cfg options

// WithConfigSources modify the default config sources for your platform.
func WithConfigSources(sources []string) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.ConfigSources = sources
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithOnPrepare provide an OnPrepare callable for the Cfg instance.
func WithOnPrepare(onPrepare func(*network.Driver) error) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.OnPrepare = onPrepare
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithDedicatedConnection set dedicated connection for Cfg instance.
func WithDedicatedConnection(dedicatedConnection bool) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.DedicatedConnection = dedicatedConnection
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithIgnoreVersion set ignore version for Cfg instance.
func WithIgnoreVersion(ignoreVersion bool) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.IgnoreVersion = ignoreVersion
			return nil
		}

		return ErrIgnoredOption
	}
}

// platform specific options

func setPlatformAttr(attrName string, attrValue, p interface{}) error {
	_, ok := p.(*Cfg)

	if ok {
		// this func only sets attrs for the platforms, so if we see a *Cfg we know we can bail
		return ErrIgnoredOption
	}

	v := reflect.ValueOf(p).Elem()

	fieldNames := map[string]int{}

	for i := 0; i < v.NumField(); i++ {
		fieldNames[v.Type().Field(i).Name] = i
	}

	attrIndex := -1

	for name, i := range fieldNames {
		if name == attrName {
			attrIndex = i
			break
		}
	}

	if attrIndex == -1 {
		// for some reason the platform doesnt have the specified attribute, this should *not*
		// happen... in theory :)
		return ErrInvalidPlatformAttr
	}

	fieldVal := v.Field(attrIndex)
	fieldVal.Set(reflect.ValueOf(attrValue))

	return nil
}

// WithVersionPattern set version pattern for the platform instance.
func WithVersionPattern(versionPattern *regexp.Regexp) Option {
	return func(p interface{}) error {
		err := setPlatformAttr("VersionPattern", versionPattern, p)

		if err != nil {
			if !errors.Is(err, ErrIgnoredOption) {
				return err
			}
		}

		return nil
	}
}

// WithFilesystem set string name of filesystem to use.
func WithFilesystem(fs string) Option {
	return func(p interface{}) error {
		err := setPlatformAttr("Filesystem", fs, p)

		if err != nil {
			if !errors.Is(err, ErrIgnoredOption) {
				return err
			}
		}

		return nil
	}
}

// WithCandidateConfigFilename set the default candidate config filename for your platform. Better
// to use only in the testing env where a definite constant filename of config file on the device
// is helpful. If, however, you use this for "normal" operations, make sure you abort any loaded
// configuration between if doing subsequent load operations as the candidate config will remain the
// same when using this option!
func WithCandidateConfigFilename(fn string) Option {
	return func(p interface{}) error {
		err := setPlatformAttr("CandidateConfigFilename", fn, p)

		if err != nil {
			if !errors.Is(err, ErrIgnoredOption) {
				return err
			}
		}

		return nil
	}
}

// operation specific options

// OperationOptions struct for options for any "operation" (LoadConfig, CommitConfig, etc.).
type OperationOptions struct {
	Source              string
	DiffColorize        bool
	DiffSideBySideWidth int
	AutoClean           bool
	Kwargs              map[string]string
}

// OperationOption function to set options for cfg operations.
type OperationOption func(*OperationOptions)

// WithConfigSource set version pattern for the platform instance.
func WithConfigSource(source string) OperationOption {
	return func(o *OperationOptions) {
		o.Source = source
	}
}

// WithDiffColorize set colorize attribute of diff response object.
func WithDiffColorize(c bool) OperationOption {
	return func(o *OperationOptions) {
		o.DiffColorize = c
	}
}

// WithDiffSideBySideWidth set side by side diff width of diff response object.
func WithDiffSideBySideWidth(i int) OperationOption {
	return func(o *OperationOptions) {
		o.DiffSideBySideWidth = i
	}
}

// WithAutoClean set auto clean for IOSXE platform.
func WithAutoClean(a bool) OperationOption {
	return func(o *OperationOptions) {
		o.AutoClean = a
	}
}

// WithKwargs option that accepts a map to act like python kwargs a little.
func WithKwargs(m map[string]string) OperationOption {
	return func(o *OperationOptions) {
		o.Kwargs = m
	}
}
