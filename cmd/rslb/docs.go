package main

import (
	"fmt"
	"io"
	"os"
)

const helpText = `Usage: %s <UPSTREAM> <UPSTREAM>...

UPSTREAM: Target URL for reverse proxy. If it can be converted to a uint16, then
          it will be converted to http://localhost:<UPSTREAM>

The RSLB_PORT environment variable must be set to the port that the load
balancer will listen on.
`

func help(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, helpText, os.Args[0])
	return
}
