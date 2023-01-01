package loadbalancer

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

type LoadBalancer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

/* Abstract LoadBalancer type */
type AbstractLoadBalancer struct {
	// Name of the load balancer
	Name   string
	URL    *url.URL
	Server *http.Server

	// Mapping between URL -> Backend
	Mapping map[*url.URL]*backend.Backend
}

func (lb *AbstractLoadBalancer) InsertBackend(serverURL string) error {
	bknd, err := backend.New(serverURL)
	if err != nil {
		return err
	}

	lb.Mapping[bknd.URL] = bknd
	log.Printf("Inserted backend for %s into %s", bknd.URL.String(), lb.Name)
	return nil
}

func (lb *AbstractLoadBalancer) InitLoadBalancer() {
	// Init the backends within the lb.Mapping
	bknds := make([]*backend.Backend, len(lb.Mapping))
	for _, bknd := range lb.Mapping {
		bknds = append(bknds, bknd)
	}
	backend.InitBackends(bknds)

	// Start the LoadBalancer itself
	lb.Server.ListenAndServe()
}
