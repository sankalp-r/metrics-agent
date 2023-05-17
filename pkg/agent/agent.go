package agent

import (
	"context"
	"github.com/sankalp-r/metrics-agent/pkg/metrics"
	"github.com/sankalp-r/metrics-agent/pkg/reporter"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Agent interface {
	Start()
}

// minimalAgent represents the simple implementation of Agent interface
type minimalAgent struct {
	samplingDuration time.Duration
	reporter         reporter.Reporter
	metricsSources   []metrics.Source
}

// Start the minimalAgent
func (m *minimalAgent) Start() {
	go m.reporter.Report()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	m.sample(ctx)
}

// sample the data from respective sources
func (m *minimalAgent) sample(ctx context.Context) {
	ticker := time.NewTicker(m.samplingDuration)
	var wg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			for _, src := range m.metricsSources {
				wg.Add(1)

				src := src
				go func() {
					defer wg.Done()
					res, err := src.Collect()
					if err != nil {
						zap.L().Error(err.Error())
						return
					}
					m.reporter.Write(res)
				}()
			}

		case <-ctx.Done():
			zap.L().Info("graceful shutting down agent...")
			wg.Wait()
			m.reporter.Close()
			return

		}
	}

}
