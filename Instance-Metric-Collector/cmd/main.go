package main

import (
	"context"
	"gpu-metric-collector/pkg/collector"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	quitChan := make(chan struct{})
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())


	metricCollector := collector.NewMetricCollector()

	wg.Add(1)
	metricCollector.RunMetricCollector(ctx, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	close(quitChan)
	cancel()
	wg.Wait()
	os.Exit(0)
}
