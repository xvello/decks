package contextkeys

import (
	"sort"
	"testing"

	"github.com/OneOfOne/xxhash"

	"github.com/DataDog/datadog-agent/pkg/metrics"
)

var xxh *xxhash.XXHash64

func init() {
	xxh = xxhash.New64()
}

func xxHash(metricSample *metrics.MetricSample) uint64 {
	xxh.Reset()

	if len(metricSample.Tags) < 15 {
		selectionSort(metricSample.Tags)
	} else {
		sort.Strings(metricSample.Tags)
	}

	xxh.WriteString(metricSample.Name)
	xxh.WriteString(",")
	for _, t := range metricSample.Tags {
		xxh.WriteString(t)
		xxh.WriteString(",")
	}
	xxh.WriteString(metricSample.Host)

	return xxh.Sum64()
}

func BenchmarkXXHash64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = xxHash(testSample)
	}
}
