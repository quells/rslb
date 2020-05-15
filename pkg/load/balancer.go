// Package load provides an HTTP reverse proxy load balancer
// that serves requests in a round robin order.
package load

import (
	"net/http"
	"net/http/httputil"
	"sync"
)

// Balancer proxies incoming requests in round robin order.
type Balancer struct {
	proxies []*httputil.ReverseProxy
	idx     int
	lock    sync.Mutex
}

// NewBalancer constructs a new Balancer across the provided proxies.
func NewBalancer(proxies []*httputil.ReverseProxy) *Balancer {
	var idx int
	var lock sync.Mutex
	lb := Balancer{proxies, idx, lock}
	return &lb
}

func (lb *Balancer) getProxy() *httputil.ReverseProxy {
	lb.lock.Lock()
	i := lb.idx
	lb.idx++
	if lb.idx >= len(lb.proxies) {
		lb.idx = 0
	}
	lb.lock.Unlock()
	return lb.proxies[i]
}

func (lb *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.getProxy().ServeHTTP(w, r)
}
