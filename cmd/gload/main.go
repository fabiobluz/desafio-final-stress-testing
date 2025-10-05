package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/acme/gload/internal/cli"
	"github.com/acme/gload/internal/report"
	"github.com/acme/gload/internal/runner"
)

func main() {
	cfg, err := cli.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("flags error: %v", err)
	}

	ctx := context.Background()
	start := time.Now()
	sum, err := runner.Run(ctx, cfg)
	if err != nil {
		log.Fatalf("run error: %v", err)
	}
	sum.TotalTime = time.Since(start)

	if err := report.Render(os.Stdout, cfg.Format, sum); err != nil {
		log.Fatalf("render error: %v", err)
	}
	fmt.Println()
}
