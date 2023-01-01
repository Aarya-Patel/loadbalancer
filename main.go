package main

import (
	"log"

	"github.com/Aarya-Patel/loadbalancer/pkg/loadbalancer"
)

var (
	lbURL = "http://localhost:8080"
	hosts = []string{
		"http://localhost:8000",
		"http://localhost:8001",
		"http://localhost:8002",
	}
)

func main() {
	rrlb, err := loadbalancer.NewRoundRobinLoadBalancer("rrlb", lbURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range hosts {
		rrlb.InsertBackend(host)
	}

	rrlb.InitLoadBalancer()
}
