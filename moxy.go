package moxy

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
)

// NewReverseProxy returns a new ReverseProxy that load-balances the proxy requests between multiple hosts
// It also allows to define a chain of filter functions to process the outgoing response(s)
func NewReverseProxy(hosts []string, filters []FilterFunc) *ReverseProxy {

	director := func(request *http.Request) {
		host, err := pick(hosts)
		if err != nil {
			log.Fatal(err)
		}
		u, err := url.Parse(host)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("httpproxy: request: %s %s%s %s; upstream %s://%s", request.Method,
			request.Host, request.RequestURI, request.Header, u.Scheme, u.Host)
		request.URL.Scheme = u.Scheme
		request.URL.Host = u.Host
		request.Host = u.Host
	}

	tlsConfig := &tls.Config{}
	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &ReverseProxy{
		Transport: transport,
		Director:  director,
		Filters:   filters,
	}
}
