package parsing

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fieldSeparator = []byte("|")
var valueSeparator = []byte(":")

type MetricSample struct {
	Name       string
	Value      float64
	RawValue   string
	Tags       []string
	Host       string
	SampleRate float64
	Timestamp  float64
}

var testMetric = []byte("daemon:666|g|#sometag1:somevalue1,sometag2:somevalue2")

func parseMetricMessage(message []byte) (*MetricSample, error) {
	// daemon:666|g|#sometag1:somevalue1,sometag2:somevalue2
	// daemon:666|g|@0.1|#sometag:somevalue"

	splitMessage := bytes.Split(message, []byte("|"))

	if len(splitMessage) < 2 || len(splitMessage) > 4 {
		return nil, fmt.Errorf("invalid field number for %q", message)
	}

	// Extract name, value and type
	rawNameAndValue := bytes.Split(splitMessage[0], []byte(":"))

	if len(rawNameAndValue) != 2 {
		return nil, fmt.Errorf("invalid field format for %q", message)
	}

	rawName, rawValue, rawType := rawNameAndValue[0], rawNameAndValue[1], splitMessage[1]
	if len(rawName) == 0 || len(rawValue) == 0 || len(rawType) == 0 {
		return nil, fmt.Errorf("invalid metric message format: empty 'name', 'value' or 'text' field")
	}

	// Metadata
	var metricTags []string
	var host string
	rawSampleRate := []byte("1")
	if len(splitMessage) > 2 {
		rawMetadataFields := splitMessage[2:]

		for i := range rawMetadataFields {
			if len(rawMetadataFields[i]) < 2 {
				continue
			}

			if bytes.HasPrefix(rawMetadataFields[i], []byte("#")) {
				metricTags, host = parseTagsStream(rawMetadataFields[i], true)
			} else if bytes.HasPrefix(rawMetadataFields[i], []byte("@")) {
				rawSampleRate = rawMetadataFields[i][1:]
			} else {
			}
		}
	}

	// Casting
	metricName := string(rawName)

	metricSampleRate, err := strconv.ParseFloat(string(rawSampleRate), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid sample value for %q", message)
	}

	sample := &MetricSample{
		Name:       metricName,
		Tags:       metricTags,
		Host:       host,
		SampleRate: metricSampleRate,
		Timestamp:  0,
	}

	return sample, nil
}

func BenchmarkParseMetricCurrent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := parseMetricMessage(testMetric)
		assert.Nil(b, err)
	}
}

func parseMetricMessageSplit2(message []byte) (*MetricSample, error) {
	// daemon:666|g|#sometag1:somevalue1,sometag2:somevalue2
	// daemon:666|g|@0.1|#sometag:somevalue"

	separatorCount := bytes.Count(message, fieldSeparator)
	if separatorCount < 2 || separatorCount > 4 {
		return nil, fmt.Errorf("invalid field number for %q", message)
	}

	// Extract name, value and type
	rawNameAndValue, remainder := split2(message, fieldSeparator)
	rawName, rawValue := split2(rawNameAndValue, valueSeparator)

	if rawValue == nil {
		return nil, fmt.Errorf("invalid field format for %q", message)
	}

	rawType, remainder := split2(remainder, fieldSeparator)
	if len(rawName) == 0 || len(rawValue) == 0 || len(rawType) == 0 {
		return nil, fmt.Errorf("invalid metric message format: empty 'name', 'value' or 'text' field")
	}

	// Metadata
	var metricTags []string
	var host string
	var rawMetadataField []byte
	sampleRate := 1.0

	for {
		rawMetadataField, remainder = split2(remainder, fieldSeparator)

		if bytes.HasPrefix(rawMetadataField, []byte("#")) {
			metricTags, host = parseTagsStream(rawMetadataField, true)
		} else if bytes.HasPrefix(rawMetadataField, []byte("@")) {
			rawSampleRate := rawMetadataField[1:]
			var err error
			sampleRate, err = strconv.ParseFloat(string(rawSampleRate), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid sample value for %q", message)
			}
		}

		if remainder == nil {
			break
		}
	}

	metricName := string(rawName)

	sample := &MetricSample{
		Name:       metricName,
		Tags:       metricTags,
		Host:       host,
		SampleRate: sampleRate,
		Timestamp:  0,
	}

	return sample, nil
}

func BenchmarkParseMetricSplit2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := parseMetricMessageSplit2(testMetric)
		assert.Nil(b, err)
	}
}
