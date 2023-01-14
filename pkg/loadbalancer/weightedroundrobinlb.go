package loadbalancer

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

type WeightedRoundRobinLoadBalancer struct {
	abstractLoadBalancer
	weights map[string]uint
	index   uint
	count   uint
	mux     sync.Mutex
}

func NewWeightedRoundRobinLoadBalancer(name string, serverURL string) (*WeightedRoundRobinLoadBalancer, error) {
	targetURL, err := url.Parse(serverURL)
	if err != nil {
		return &WeightedRoundRobinLoadBalancer{}, err
	}

	lb := WeightedRoundRobinLoadBalancer{
		abstractLoadBalancer: abstractLoadBalancer{},
		weights:              make(map[string]uint),
		index:                0,
		count:                0,
		mux:                  sync.Mutex{},
	}

	abstractLB := abstractLoadBalancer{
		Name: name,
		URL:  targetURL,
		Server: &http.Server{
			Addr:    targetURL.Host,
			Handler: http.HandlerFunc(lb.GenerateLBServeHTTP()),
		},
		ServerURLs: []string{},
		Mapping:    make(map[string]*backend.Backend),
	}

	// Set the abstractLB after we GenerateLBServeHTTP
	lb.abstractLoadBalancer = abstractLB

	return &lb, nil
}

func (lb *WeightedRoundRobinLoadBalancer) GenerateLBServeHTTP() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if len(lb.Mapping) == 0 {
			io.WriteString(rw, "There are no backends available to process your request!")
			return
		}

		currentServerURL := lb.ServerURLs[lb.index]

		// Update the parameters for the next packet
		lb.mux.Lock()
		lb.count = (lb.count + 1) % lb.weights[currentServerURL]
		if lb.count == 0 {
			lb.index = (lb.index + 1) % uint(len(lb.Mapping))
			lb.count = 0
		}
		lb.mux.Unlock()

		bknd := lb.Mapping[currentServerURL]
		if bknd.IsHealthy() {
			bknd.ReverseProxy.ServeHTTP(rw, req)
			return
		}
		io.WriteString(rw, "There are no healthy backends available to process your request!")
	}
}

func (lb *WeightedRoundRobinLoadBalancer) InsertBackend(serverURL string, weight uint) error {
	if weight <= 0 {
		return errors.New("Weight must be positive")
	}

	err := lb.abstractLoadBalancer.InsertBackend(serverURL)
	if err != nil {
		return err
	}

	targetURL, err := url.Parse(serverURL)
	if err != nil {
		return err
	}

	lb.weights[targetURL.String()] = weight
	return nil
}
