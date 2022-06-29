package github

import (
	"net"
	"net/http"
	"time"
)

func newTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConns:          0,
		IdleConnTimeout:       15 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
	}
}

func newClient(transport *http.Transport) *http.Client {
	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
}
