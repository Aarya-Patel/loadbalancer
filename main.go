package main

import (
	"log"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

var (
	hosts = []string{
		"http://localhost:8000",
		"http://localhost:8001",
		"http://localhost:8002",
	}
)

func main() {
	backends := []*backend.Backend{}

	for _, host := range hosts {
		backend, err := backend.New(host)
		if err != nil {
			log.Fatal(err)
		}

		backends = append(backends, backend)
	}

	backend.InitBackends(backends)
	for {
	}
}
