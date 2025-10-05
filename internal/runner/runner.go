package runner

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/acme/gload/internal/httpc"
	"github.com/acme/gload/internal/report"
	"github.com/acme/gload/internal/worker"
)

type Config struct {
	URL         string
	Requests    int
	Concurrency int
	Timeout     time.Duration
	Method      string
	Headers     http.Header
	Body        []byte
	Format      string
}

func Run(ctx context.Context, cfg Config) (report.Summary, error) {
	jobs := make(chan int, cfg.Requests)
	results := make(chan worker.Result, cfg.Requests)

	client := httpc.New(cfg.Timeout)
	spec := worker.RequestSpec{
		Method:  cfg.Method,
		URL:     cfg.URL,
		Headers: cfg.Headers,
		Body:    cfg.Body,
	}

	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker.Loop(ctx, id, jobs, results, client, spec)
		}(i)
	}

	for i := 0; i < cfg.Requests; i++ { jobs <- i }
	close(jobs)
	wg.Wait()
	close(results)

	var acc []report.Result
	for r := range results {
		acc = append(acc, report.Result{Status: r.Status, Latency: r.Latency})
	}
	return report.Collect(acc), nil
}
