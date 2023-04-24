package main

import (
	"github.com/pborman/getopt/v2"
	"io"
)

type config struct {
	verboseFlag         *bool
	helpFlag            *bool
	versionFlag         *bool
	followRedirectsFlag *bool

	set *getopt.Set
}

func newConfig() *config {
	c := &config{}

	c.set = getopt.New()

	c.verboseFlag = c.set.BoolLong("verbose", 'v', "Make the operation more talkative")
	c.helpFlag = c.set.BoolLong("help", 'h', "Show this help")
	c.versionFlag = c.set.BoolLong("version", 'V', "Show version information and quit")
	c.followRedirectsFlag = c.set.BoolLong("location", 'L', "Follow redirects")

	return c
}

func (c *config) parse(args []string) {
	c.set.Parse(args)
}

func (c *config) usage(w io.Writer) {
	c.set.PrintUsage(w)
}
