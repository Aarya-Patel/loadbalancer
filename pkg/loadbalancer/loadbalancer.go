package loadbalancer

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

type LoadBalancer interface {
	GenerateLBServeHTTP() func(http.ResponseWriter, *http.Request)
}

/* Abstract LoadBalancer type */
type abstractLoadBalancer struct {
	// Name of the load balancer (for logs)
	Name   string
	URL    *url.URL
	Server *http.Server

	// Mapping between URL -> Backend
	Mapping map[*url.URL]*backend.Backend
}

func (lb *abstractLoadBalancer) InsertBackend(serverURL string) error {
	bknd, err := backend.New(serverURL)
	if err != nil {
		return err
	}

	lb.Mapping[bknd.URL] = bknd
	log.Printf("Inserted backend on %s into %s", bknd.URL.String(), lb.Name)
	return nil
}

func (lb *abstractLoadBalancer) InitLoadBalancer() {
	// Init the backends within the lb.Mapping
	bknds := []*backend.Backend{}
	for _, bknd := range lb.Mapping {
		bknds = append(bknds, bknd)
	}
	backend.InitBackends(bknds)

	// Start the LoadBalancer itself
	log.Printf("Load balancer named '%s' is listening on %s", lb.Name, lb.URL.Host)
	lb.Server.ListenAndServe()
}
