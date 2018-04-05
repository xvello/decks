package contextkeys

import (
	"bytes"
	"sort"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/metrics"
)

///// Current DSD6 hash

func generateContextKey(metricSample *metrics.MetricSample) string {
	// Pre-compute the size of the buffer we'll need, and allocate a buffer of that size
	bufferSize := len(metricSample.Name) + 1
	for k := range metricSample.Tags {
		bufferSize += len(metricSample.Tags[k]) + 1
	}
	bufferSize += len(metricSample.Host)
	buffer := bytes.NewBuffer(make([]byte, 0, bufferSize))

	sort.Strings(metricSample.Tags)
	// write the context items to the buffer, and return it as a string
	buffer.WriteString(metricSample.Name)
	buffer.WriteString(",")
	for k := range metricSample.Tags {
		buffer.WriteString(metricSample.Tags[k])
		buffer.WriteString(",")
	}
	buffer.WriteString(metricSample.Host)

	return buffer.String()
}

func BenchmarkCurrent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = generateContextKey(testSample)
	}
}
