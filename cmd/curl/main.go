package main

import (
	"context"
	"fmt"
	"github.com/lsmoura/curl"
	"net/http"
	"os"
	"runtime/debug"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "%s: try '%s --help' for more information", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
			fmt.Fprintf(os.Stderr, "Stack Trace:\n")
			os.Stderr.Write(debug.Stack())
			fmt.Fprintf(os.Stderr, "\n")
			os.Exit(1)
		}
	}()

	ctx := context.Background()

	conf := newConfig()
	conf.parse(os.Args)

	if *conf.helpFlag {
		conf.usage(os.Stdout)
		os.Exit(0)
	}
	if *conf.versionFlag {
		fmt.Printf("%s %s (%s)\n", os.Args[0], version, buildDate)
		os.Exit(0)
	}

	for _, path := range os.Args[1:] {
		req := curl.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if *conf.verboseFlag {
			req.Verbose = true
		}
		if *conf.followRedirectsFlag {
			req.FollowRedirects = true
		}

		data, err := curl.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s", os.Args[0], err)
			os.Exit(1)
		}

		out := req.Out
		if out == nil {
			out = os.Stdout
		}
		fmt.Fprintf(out, "%s", data)

		fmt.Println()
	}
}
