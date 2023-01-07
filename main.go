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
	hcTimeout  = 1e9
	hcInterval = 3e9
)

func main() {
	rrlb, err := loadbalancer.NewRoundRobinLoadBalancer("rrlb", lbURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range hosts {
		rrlb.InsertBackend(host)
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

	go hc.PingBackends()
	rrlb.InitLoadBalancer()
}
