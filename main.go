package main

import (
	"log"
	"time"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
	"github.com/Aarya-Patel/loadbalancer/pkg/healthcheck"
	"github.com/Aarya-Patel/loadbalancer/pkg/loadbalancer"
)

var (
	lbURL = "http://localhost:8080"
	hosts = []string{
		"http://localhost:8000",
		"http://localhost:8001",
		"http://localhost:8002",
	}
	hcTimeout  = 1 * time.Second
	hcInterval = 2 * time.Second
)

func main() {
	// rrlb, err := loadbalancer.NewRoundRobinLoadBalancer("rrlb", lbURL)
	rrlb, err := loadbalancer.NewRandomizedLoadBalancer("rrlb", lbURL)
	// rrlb, err := loadbalancer.NewWeightedRoundRobinLoadBalancer("wrrlb", lbURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range hosts {
		err := rrlb.InsertBackend(host)
		// err := rrlb.InsertBackend(host, uint(i+1))
		if err != nil {
			log.Fatal(err)
		}
	}

	allBackends := []*backend.Backend{}
	for _, bknd := range rrlb.Mapping {
		allBackends = append(allBackends, bknd)
	}

	hc, err := healthcheck.New(
		healthcheck.WithBackends(allBackends),
		healthcheck.WithInterval(time.Duration(hcInterval)),
		healthcheck.WithTimeout(time.Duration(hcTimeout)),
	)
	if err != nil {
		log.Fatal(err)
	}

	hc.StartHealthChecks()
	rrlb.InitLoadBalancer()
}
