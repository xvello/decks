package contextkeys

import (
	"github.com/DataDog/datadog-agent/pkg/metrics"
)

var testSample *metrics.MetricSample

var separator = ","

func init() {
	testSample = &metrics.MetricSample{
		Name:  "docker.mem.rss",
		Value: 6.602752e+06,
		Mtype: metrics.GaugeType,
		Tags: []string{"container_id:4f5c7cdb401e9a7868a20e814d257e1fa69113bc4d5a285368141f59dc18643a",
			"container_name:keen_lamarr",
			"docker_image:redis:latest",
			"image_name:redis",
			"image_tag:latest",
		},
		Host:       "ci-xaviervello",
		SampleRate: 0,
		Timestamp:  1513089025,
	}
}

func bubbleSortStrings(arr []string) {
	var tmp string
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr)-1; j++ {
			if arr[j] > arr[j+1] {
				tmp = arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = tmp
			}
		}
	}
}

func Swap(arrayzor []string, i, j int) {
	tmp := arrayzor[j]
	arrayzor[j] = arrayzor[i]
	arrayzor[i] = tmp
}

func selectionSort(array []string) {
	for i := 0; i < len(array)-1; i++ {
		min := i
		for j := i + 1; j < len(array)-1; j++ {
			if array[j] < array[min] {
				min = j
			}
		}
		Swap(array, i, min)
	}
}
