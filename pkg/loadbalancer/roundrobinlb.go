package loadbalancer

import "net/http"

type RoundRobinLoadBalancer struct {
	AbstractLoadBalancer
	index int
}

func NewRoundRobinLoadBalancer(serverURL string) (*RoundRobinLoadBalancer, error) {
	return &RoundRobinLoadBalancer{}, nil
}

func (rrlb *RoundRobinLoadBalancer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// TODO
}
