package healthcheck

import (
	"time"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

type Option func(*HealthCheck)

type HealthCheck struct {
	Backends []*backend.Backend
	Timeout  time.Duration
	Interval time.Duration
}
