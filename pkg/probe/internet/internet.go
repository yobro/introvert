package internet

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/yobro/introvert/pkg/probe"
)

var (
	latency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "introvert_internet_latency_seconds",
		Help:    "internet request latency seconds",
		Buckets: prometheus.ExponentialBuckets(0.02, 1.4, 15),
	})

	requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "introvert_internet_requests_total",
		Help: "requests total",
	})

	requestFailuresTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "introvert_internet_request_failures_total",
		Help: "total request failures",
	})
)

const (
	defaultURL      = "http://google.com"
	defaultTimeout  = 5 * time.Second
	defaultInterval = 15 * time.Second
)

// Probe probes internet provider service
type Probe struct {
	URL      string
	Interval time.Duration
	Timeout  time.Duration
}

// DefaultProbe returns default probe
func DefaultProbe() probe.Probe {
	return &Probe{
		URL:      defaultURL,
		Timeout:  defaultTimeout,
		Interval: defaultInterval,
	}
}

// Run implements Probe interface
func (p *Probe) Run(ctx context.Context) error {

	if _, err := url.Parse(p.URL); err != nil {
		return err
	}

	ticker := time.NewTicker(p.Interval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			requestsTotal.Inc()
			timer := prometheus.NewTimer(latency)
			if err := p.ping(ctx); err != nil {
				requestFailuresTotal.Inc()
				log.Printf("error pinging %s: %v", p.URL, err)
			}
			timer.ObserveDuration()
		}
	}
}

func (p *Probe) ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Warning: unable to read message body: %v", err)
		}
		return fmt.Errorf("http status code %d: %s", resp.StatusCode, string(b))
	}

	return nil
}
