package healthcheck

import (
	"log"
	"net/http"
	"time"

	"github.com/Aarya-Patel/loadbalancer/internal/backend"
)

func New(options ...Option) (*HealthCheck, error) {
	hc := HealthCheck{
		Backends: []*backend.Backend{},
		Timeout:  0.0,
	}

	for _, option := range options {
		option(&hc)
	}

	return &hc, nil
}

func WithTimeout(timeout time.Duration) Option {
	return func(hc *HealthCheck) {
		hc.Timeout = timeout
	}
}

func WithInterval(interval time.Duration) Option {
	return func(hc *HealthCheck) {
		hc.Interval = interval
	}
}

func WithBackends(bknds []*backend.Backend) Option {
	return func(hc *HealthCheck) {
		hc.Backends = bknds
	}
}

func (hc *HealthCheck) PingBackends() {
	doneChannels := make([]chan backend.Status, len(hc.Backends))
	for idx := range doneChannels {
		doneChannels[idx] = make(chan backend.Status)
	}

	for idx, bknd := range hc.Backends {
		go hc.pingBackend(bknd, &doneChannels[idx])

	waitForPingBackend:
		for {
			select {
			case status := <-(doneChannels[idx]):
				log.Printf("Ping to %s returned with backend status of %s", bknd.URL.Host, status.String())
				if bknd.Status != status {
					bknd.UpdateBackendStatus(status)
				}
				break waitForPingBackend
			}
		}
	}
}

func (hc *HealthCheck) pingBackend(bknd *backend.Backend, done *chan backend.Status) {
	time.AfterFunc(hc.Timeout, func() {
		*done <- backend.Unhealthy
	})

	resp, err := http.Get(bknd.URL.String())
	if err != nil {
		*done <- backend.Unhealthy
		return
	}

	if resp.StatusCode == http.StatusOK {
		*done <- backend.Healthy
	} else {
		*done <- backend.Unhealthy
	}

}
