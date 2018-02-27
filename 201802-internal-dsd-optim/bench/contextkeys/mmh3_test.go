package contextkeys

import (
	"encoding/binary"
	"io"
	"sort"
	"testing"

	"github.com/DataDog/mmh3"

	"github.com/DataDog/datadog-agent/pkg/metrics"
)

var mmh *mmh3.HashWriter128

func init() {
	mmh = new(mmh3.HashWriter128)
}

////// dd-go mmh3 hash

// Handle is a hash of all the inputs to a context. Consider it a context's
// natural key crammed into 20 bytes.
const ContextHandleLength = 20

type Handle [ContextHandleLength]byte

// NewHandle will create the datadog context handle for the given inputs.
func NewHandle(metricSample *metrics.MetricSample) Handle {
	mmh := mmh3.HashWriter128{}
	b := &mmh

	sort.Strings(metricSample.Tags)

	io.WriteString(b, metricSample.Name)
	io.WriteString(b, ",")
	for _, t := range metricSample.Tags {
		io.WriteString(b, t)
		io.WriteString(b, ",")
	}
	io.WriteString(b, metricSample.Host)

	var h Handle
	binary.LittleEndian.PutUint32(h[0:], 0)
	mmh.Sum(h[4:4])
	return h
}

func BenchmarkHandle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewHandle(testSample)
	}
}

// NewHandle will create the datadog context handle for the given inputs.
func NewHandleReuse(metricSample *metrics.MetricSample) Handle {
	mmh.Reset()

	sort.Strings(metricSample.Tags)

	io.WriteString(mmh, metricSample.Name)
	io.WriteString(mmh, ",")
	for _, t := range metricSample.Tags {
		io.WriteString(mmh, t)
		io.WriteString(mmh, ",")
	}
	io.WriteString(mmh, metricSample.Host)

	var h Handle
	binary.LittleEndian.PutUint32(h[0:], 0)
	mmh.Sum(h[4:4])
	return h
}

func BenchmarkHandleReuse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewHandleReuse(testSample)
	}
}

// NewHandle will create the datadog context handle for the given inputs.
func NewHandleReuseNoIface(metricSample *metrics.MetricSample) [16]byte {
	mmh.Reset()

	sort.Strings(metricSample.Tags)

	mmh.WriteString(metricSample.Name)
	mmh.WriteString(",")
	for _, t := range metricSample.Tags {
		mmh.WriteString(t)
		mmh.WriteString(",")
	}
	mmh.WriteString(metricSample.Host)

	//binary.LittleEndian.PutUint32(h[0:], 0)
	var hash [16]byte
	mmh.Sum(hash[0:0])
	return hash
}

func BenchmarkHandleReuseNoIface(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewHandleReuseNoIface(testSample)
	}
}

// NewHandle will create the datadog context handle for the given inputs.
func NewHandleReuseNoIfaceBubble(metricSample *metrics.MetricSample) [16]byte {
	mmh.Reset()

	if len(metricSample.Tags) < 15 {
		bubbleSortStrings(metricSample.Tags)
	} else {
		sort.Strings(metricSample.Tags)
	}

	mmh.WriteString(metricSample.Name)
	mmh.WriteString(",")
	for _, t := range metricSample.Tags {
		mmh.WriteString(t)
		mmh.WriteString(",")
	}
	mmh.WriteString(metricSample.Host)

	//binary.LittleEndian.PutUint32(h[0:], 0)
	var hash [16]byte
	mmh.Sum(hash[0:0])
	return hash
}

func BenchmarkHandleReuseNoIfaceBubble(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewHandleReuseNoIfaceBubble(testSample)
	}
}

func NewHandleReuseNoIfaceSelection(metricSample *metrics.MetricSample) [16]byte {
	mmh.Reset()

	if len(metricSample.Tags) < 15 {
		selectionSort(metricSample.Tags)
	} else {
		sort.Strings(metricSample.Tags)
	}

	mmh.WriteString(metricSample.Name)
	mmh.WriteString(",")
	for _, t := range metricSample.Tags {
		mmh.WriteString(t)
		mmh.WriteString(",")
	}
	mmh.WriteString(metricSample.Host)

	//binary.LittleEndian.PutUint32(h[0:], 0)
	var hash [16]byte
	mmh.Sum(hash[0:0])
	return hash
}

func BenchmarkHandleReuseNoIfaceSelection(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = NewHandleReuseNoIfaceSelection(testSample)
	}
}
