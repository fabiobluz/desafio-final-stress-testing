package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

type Result struct {
	Status  int
	Latency time.Duration
}

type Summary struct {
	TotalRequests int           `json:"total_requests"`
	Success200    int           `json:"success_200"`
	StatusDist    map[int]int   `json:"status_distribution"`
	TotalTime     time.Duration `json:"-"`
	AvgLatency    time.Duration `json:"-"`
	P95           time.Duration `json:"-"`
	P99           time.Duration `json:"-"`
}

// JSONSummary é usado apenas para serialização JSON com tempos legíveis
type JSONSummary struct {
	TotalRequests int         `json:"total_requests"`
	Success200    int         `json:"success_200"`
	StatusDist    map[int]int `json:"status_distribution"`
	TotalTime     string      `json:"total_time"`
	AvgLatency    string      `json:"avg_latency,omitempty"`
	P95           string      `json:"p95_latency,omitempty"`
	P99           string      `json:"p99_latency,omitempty"`
}

func Collect(results []Result) Summary {
	s := Summary{StatusDist: make(map[int]int)}
	var lats []time.Duration
	for _, r := range results {
		s.TotalRequests++
		if r.Status == 200 {
			s.Success200++
		}
		s.StatusDist[r.Status]++
		lats = append(lats, r.Latency)
	}
	if len(lats) > 0 {
		var sum time.Duration
		for _, d := range lats {
			sum += d
		}
		s.AvgLatency = sum / time.Duration(len(lats))
		s.P95 = percentile(lats, 95)
		s.P99 = percentile(lats, 99)
	}
	return s
}

func Render(w io.Writer, format string, s Summary) error {
	switch format {
	case "json":
		jsonSummary := JSONSummary{
			TotalRequests: s.TotalRequests,
			Success200:    s.Success200,
			StatusDist:    s.StatusDist,
			TotalTime:     s.TotalTime.String(),
			AvgLatency:    s.AvgLatency.String(),
			P95:           s.P95.String(),
			P99:           s.P99.String(),
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(jsonSummary)
	default:
		_, err := fmt.Fprintf(w,
			"\nTotal time: %v\nTotal requests: %d\nHTTP 200: %d\nStatus distribution: %v\nAvg latency: %v\nP95: %v\nP99: %v\n",
			s.TotalTime, s.TotalRequests, s.Success200, s.StatusDist, s.AvgLatency, s.P95, s.P99,
		)
		return err
	}
}

func percentile(ds []time.Duration, p int) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	sort.Slice(ds, func(i, j int) bool { return ds[i] < ds[j] })
	idx := (len(ds)*p+99)/100 - 1 // ceil
	if idx < 0 {
		idx = 0
	}
	if idx >= len(ds) {
		idx = len(ds) - 1
	}
	return ds[idx]
}
