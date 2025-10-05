package report

import (
	"testing"
	"time"
)

func TestCollectAndPercentiles(t *testing.T) {
	results := []Result{
		{Status: 200, Latency: 10 * time.Millisecond},
		{Status: 200, Latency: 20 * time.Millisecond},
		{Status: 500, Latency: 30 * time.Millisecond},
		{Status: 404, Latency: 40 * time.Millisecond},
		{Status: 200, Latency: 50 * time.Millisecond},
	}

	s := Collect(results)

	if s.TotalRequests != 5 {
		t.Fatalf("TotalRequests = %d, want 5", s.TotalRequests)
	}
	if s.Success200 != 3 {
		t.Fatalf("Success200 = %d, want 3", s.Success200)
	}
	if s.StatusDist[500] != 1 || s.StatusDist[404] != 1 {
		t.Fatalf("StatusDist unexpected: %#v", s.StatusDist)
	}
	if s.AvgLatency <= 0 || s.P95 <= 0 || s.P99 <= 0 {
		t.Fatalf("latency stats should be > 0, got avg=%v p95=%v p99=%v", s.AvgLatency, s.P95, s.P99)
	}
}
