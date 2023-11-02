package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gorilla/mux"

	"github.com/AhmedSalah101/tareek/util"
)

var pl = fmt.Println

type Target struct {
	addr   string
	router *mux.Router
	proxy  *httputil.ReverseProxy
}

type LoadBalancer struct {
	lock            sync.RWMutex
	targets         []*Target
	roundRobinCount int
	router          *mux.Router
}

func newTarget(addr string, r *mux.Router) *Target {
	targetURL, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}

	prox := httputil.NewSingleHostReverseProxy(targetURL)
	prox.Director = func(r *http.Request) {
		r.Host = targetURL.Host
		targetQuery := targetURL.RawQuery
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		r.URL.Path, r.URL.RawPath = targetURL.Path, targetURL.RawPath
		if targetQuery == "" || r.URL.RawQuery == "" {
			r.URL.RawQuery = targetQuery + r.URL.RawQuery
		} else {
			r.URL.RawQuery = targetQuery + "&" + r.URL.RawQuery
		}
		pl("path : ", targetURL.RawPath)
	}
	rp := &Target{addr: addr, router: r, proxy: prox}
	return rp
}

func (trgt *Target) IsAlive() bool {
	return true
}

func (lb *LoadBalancer) getNextAvailableServer() *Target {
	pl("get next")
	target := lb.targets[lb.roundRobinCount%len(lb.targets)]
	lb.lock.Lock()
	for !target.IsAlive() {
		pl("inside loop")
		lb.roundRobinCount++
		target = lb.targets[lb.roundRobinCount%len(lb.targets)]

	}
	pl("out loop")
	lb.roundRobinCount++
	lb.lock.Unlock()
	pl(target)
	return target
}

type Tareek struct {
	proxies []*Target
	lbs     []*LoadBalancer
}

func (tr *Tareek) redirector(
	w http.ResponseWriter,
	r *http.Request,
) {
	for _, rp := range tr.proxies {
		matchedRoute := &mux.RouteMatch{}
		if rp.router.Match(r, matchedRoute) {
			// pl(r.URL.Path)
			rp.proxy.ServeHTTP(w, r)
			return
		}
	}
	for _, lb := range tr.lbs {
		matchedRoute := &mux.RouteMatch{}
		if lb.router.Match(r, matchedRoute) {
			// pl(r.URL.Path)
			lb.getNextAvailableServer().proxy.ServeHTTP(
				w,
				r,
			)
		}
	}
}

func (tr *Tareek) addTarget(addr string, r *mux.Router) {
	tr.proxies = append(tr.proxies, newTarget(addr, r))
}

func (tr *Tareek) addLoadBalancer(
	addrs []string,
	r *mux.Router,
) {
	tr.lbs = append(tr.lbs, newLoadBalancer(addrs, r))
}

func newLoadBalancer(
	addrs []string,
	r *mux.Router,
) *LoadBalancer {
	var targets []*Target
	for _, addr := range addrs {
		targets = append(targets, newTarget(addr, r))
	}

	return &LoadBalancer{
		targets:         targets,
		roundRobinCount: 0,
		router:          r,
	}
}

func (tr *Tareek) Start() {
	pl("Starting server....")
	http.HandleFunc("/", tr.redirector)
	http.ListenAndServe(":4000", nil)
}

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var tr Tareek
	for _, rp := range config.RPs {
		r := mux.NewRouter()
		r.Host(rp.Host).PathPrefix(rp.Endpoint)
		tr.addTarget(rp.Target, r)
	}

	for _, lb := range config.LBs {
		r := mux.NewRouter()
		r.Host(lb.Host).PathPrefix(lb.Endpoint)
		tr.addLoadBalancer(lb.Targets, r)
	}

	tr.Start()
}
