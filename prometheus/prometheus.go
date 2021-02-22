package prometheus

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/tsdb"
)

// Lite is Prometheus with only local db and web
type Lite struct {
	storage     *tsdb.DB
	queryEngine *promql.Engine
	lis         net.Listener
	mtx         sync.Mutex
	ready       chan struct{}
}

// NewPrometheusLite returns a new Prometheus lite
func NewPrometheusLite(ctx context.Context, dir string) (*Lite, error) {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}

	opts := tsdb.DefaultOptions()
	opts.RetentionDuration = int64(365 * 24 * time.Hour / time.Millisecond)
	db, err := tsdb.Open(path.Join(dir, dir), nil, prometheus.DefaultRegisterer, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric database: %v", err)
	}

	queryEngine := promql.NewEngine(promql.EngineOpts{Reg: prometheus.DefaultRegisterer})

	return &Lite{storage: db, queryEngine: queryEngine, ready: make(chan struct{})}, nil
}

// Run starts prometheus server
func (p *Lite) Run(ctx context.Context, addr string) error {
	p.mtx.Lock()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	p.lis = lis
	p.mtx.Unlock()

	ch := make(chan error)
	go func() {
		ch <- http.Serve(lis, p.Handler())
	}()

	close(p.ready)
	defer func() {
		p.ready = make(chan struct{})
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}

// Addr returns the address of the current listener
func (p *Lite) Addr() string {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.lis == nil {
		return ""
	}

	addr := p.lis.Addr().String()
	return fmt.Sprintf("http://localhost:%s", addr[strings.LastIndex(addr, ":")+1:])
}

// Ready blocks until the server is ready
func (p *Lite) Ready() <-chan struct{} {
	return p.ready
}

// QueryRange returns query between x and y time
func (p *Lite) QueryRange(ctx context.Context, promql string, start, end time.Time, step time.Duration) (interface{}, error) {

	query, err := p.queryEngine.NewRangeQuery(p.storage, promql, start, end, step)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	res := query.Exec(ctx)
	if res.Err != nil {
		return nil, err
	}

	return res.Value, nil
}
