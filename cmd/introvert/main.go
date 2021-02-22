package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/yobro/introvert/pkg/probe"
	"github.com/yobro/introvert/pkg/probe/internet"
	"github.com/yobro/introvert/pkg/server"
	"github.com/yobro/introvert/prometheus"
)

var (
	flgStorageDir        string
	flgPort              int
	flgPrometheusAddress string
)

const prometheusDir = "prometheus"

func main() {
	cmd := &cobra.Command{
		Use:   "introvert FLAGS",
		Short: "Tool for monitoring your productivity.",
		RunE:  run,
	}

	cmd.Flags().StringVar(&flgStorageDir, "storage.dir", "data", "path to store data")
	cmd.Flags().IntVar(&flgPort, "port", 9090, "server port")
	cmd.Flags().StringVar(&flgPrometheusAddress, "prometheus.address", "", "address to Prometheus where metrics are stored. If not set, a local Prometheus DB is used internally.")

	if err := cmd.Execute(); err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := os.MkdirAll(flgStorageDir, 0777); err != nil {
		return err
	}

	if flgPrometheusAddress == "" {

		promlite, err := prometheus.NewPrometheusLite(ctx, path.Join(flgStorageDir, prometheusDir))
		if err != nil {
			return fmt.Errorf("error setting up Prometheus: %v", err)
		}

		go func() {
			if err := promlite.Run(ctx, ":0"); err != nil {
				log.Fatalf("error running internal Prometheus: %v", err)
			}
		}()

		<-promlite.Ready()
		flgPrometheusAddress = promlite.Addr()
		log.Printf("internal Prometheus server running on %s", flgPrometheusAddress)
	}

	s, err := server.New(flgPort, flgPrometheusAddress)
	if err != nil {
		return err
	}

	var g sync.WaitGroup

	g.Add(1)
	go func() {
		defer g.Done()
		if err := s.Start(); err != nil {
			log.Fatalf("fatal server error: %v", err)
		}
	}()
	defer s.Stop()

	probes := []probe.Probe{
		internet.DefaultProbe(),
	}

	g.Add(len(probes))
	for _, p := range probes {
		go func(p probe.Probe) {
			defer g.Done()
			p.Run(ctx)
		}(p)
	}

	g.Wait()
	return nil
}
