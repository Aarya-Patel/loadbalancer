package loadbalancer

type WeightedRoundRobinLoadBalancer struct {
	abstractLoadBalancer
	weights []uint
}
