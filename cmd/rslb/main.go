package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"github.com/quells/rslb/pkg/load"
	"strconv"
)

func main() {
	var err error

	if len(os.Args) == 1 {
		err = help(os.Stderr)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	h := os.Args[1]
	if h == "-h" || h == "-help" || h == "--help" {
		help(os.Stdout)
		os.Exit(0)
	}

	addr := ":" + os.Getenv("RSLB_PORT")
	if addr == ":" {
		fmt.Fprintf(os.Stderr, "must set RSLB_PORT environment variable\n")
		os.Exit(2)
	}

	upstreamArgs := os.Args[1:]
	upstreamUrls := make([]*url.URL, len(upstreamArgs))
	for i, target := range upstreamArgs {
		var u *url.URL
		u, err = parseTarget(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid upstream: %v\n", err)
			os.Exit(3)
		}

		upstreamUrls[i] = u
	}

	proxies := make([]*httputil.ReverseProxy, len(upstreamUrls))
	for i, u := range upstreamUrls {
		proxies[i] = httputil.NewSingleHostReverseProxy(u)
	}

	server := http.Server{
		Addr:    addr,
		Handler: load.NewBalancer(proxies),
	}
	fmt.Fprintf(os.Stdout, "listening on %s\n", addr)
	if err = server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "http server: %v\n", err)
	}
}

func parsePort(s string) (p uint16, err error) {
	var i int64
	i, err = strconv.ParseInt(s, 10, 16)
	if err != nil {
		return
	}

	p = uint16(i)
	return
}

func parseTarget(s string) (u *url.URL, err error) {
	var p uint16
	p, err = parsePort(s)
	if err == nil {
		s = fmt.Sprintf("http://localhost:%d", p)
	}

	u, err = url.Parse(s)
	return
}
