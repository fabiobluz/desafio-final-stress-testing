package worker

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type RequestSpec struct {
	Method  string
	URL     string
	Headers http.Header
	Body    []byte
}

type Result struct {
	Status  int
	Latency time.Duration
}

type Doer interface { Do(*http.Request) (*http.Response, error) }

func Loop(ctx context.Context, id int, jobs <-chan int, out chan<- Result, client Doer, spec RequestSpec) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-jobs:
			if !ok { return }
			start := time.Now()
			req, _ := http.NewRequestWithContext(ctx, spec.Method, spec.URL, bytes.NewReader(spec.Body))
			for k, vs := range spec.Headers {
				for _, v := range vs { req.Header.Add(k, v) }
			}
			resp, err := client.Do(req)
			if err != nil {
				out <- Result{Status: 0, Latency: time.Since(start)}
				continue
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			out <- Result{Status: resp.StatusCode, Latency: time.Since(start)}
		}
	}
}
