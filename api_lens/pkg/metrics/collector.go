package metrics

import (
	"context"
	"sync"

	"api-lens/pkg/config"
	"api-lens/pkg/request"
	metricstypes "api-lens/pkg/types"
)

func CollectMetrics(ctx context.Context, config config.RequestConfig) metricstypes.MetricsCollection {
	var metrics []metricstypes.RequestMetrics
	var mu sync.Mutex
	var wg sync.WaitGroup

	batches := (config.RequestCount + config.BatchSize - 1) / config.BatchSize
	for i := 0; i < batches; i++ {
		batchSize := config.BatchSize
		if i == batches-1 {
			batchSize = config.RequestCount % config.BatchSize
			if batchSize == 0 {
				batchSize = config.BatchSize
			}
		}

		batchMetrics := make([]metricstypes.RequestMetrics, batchSize)
		for j := 0; j < batchSize; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				metric := request.SendRequest(ctx, config)

				mu.Lock()
				batchMetrics[idx] = metric
				mu.Unlock()
			}(j)
		}
		wg.Wait()

		metrics = append(metrics, batchMetrics...)
	}

	return metricstypes.MetricsCollection{
		Metrics:      metrics,
		RequestCount: config.RequestCount,
	}
}
