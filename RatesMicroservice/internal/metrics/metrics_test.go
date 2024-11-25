package metrics

import (
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestInitMetrics(t *testing.T) {
	_ = testutil.CollectAndCount(EndpointMetrics)
	InitMetrics()

	_ = testutil.CollectAndCount(EndpointMetrics)

	if count := EndpointMetrics.WithLabelValues("test_endpoint"); count != nil {
		count.Inc()
	} else {
		t.Fatal("expected a valid counter for 'test_endpoint'")
	}

	_ = testutil.CollectAndCount(EndpointMetrics)
}

func TestInitMetrics_Once(t *testing.T) {
	InitMetrics()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			InitMetrics()
		}()
	}
	wg.Wait()
}
