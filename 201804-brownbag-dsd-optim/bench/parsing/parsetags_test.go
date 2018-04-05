package parsing

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseTags(rawTags []byte, extractHost bool) ([]string, string) {
	var host string
	tags := bytes.Split(rawTags[1:], []byte(","))
	tagsList := make([]string, 0, len(tags))

	for _, tag := range tags {
		if extractHost && bytes.HasPrefix(tag, []byte("host:")) {
			host = string(tag[5:])
		} else {
			tagsList = append(tagsList, string(tag))
		}
	}
	return tagsList, host
}

func BenchmarkParseTagsCurrent(b *testing.B) {
	rawTags := []byte("backend:web,team:kenafeh,bleh")
	for n := 0; n < b.N; n++ {
		_, _ = parseTags(rawTags, true)
	}
}

var tagSeparator = []byte(",")

func split2(slice, sep []byte) ([]byte, []byte) {
	sepIndex := bytes.Index(slice, sep)
	if sepIndex == -1 {
		return slice, nil
	} else {
		return slice[:sepIndex], slice[sepIndex+1:]
	}
}

func parseTagsStream(rawTags []byte, extractHost bool) ([]string, string) {
	var host string
	var tag []byte
	tagsList := make([]string, 0, bytes.Count(rawTags, tagSeparator)+1)
	remainder := rawTags

	for {
		tag, remainder = split2(remainder, tagSeparator)
		if extractHost && bytes.HasPrefix(tag, []byte("host:")) {
			host = string(tag[5:])
		} else {
			tagsList = append(tagsList, string(tag))
		}

		if len(remainder) == 0 {
			break
		}
	}
	return tagsList, host
}

func TestParseTagsStream(t *testing.T) {
	rawTags := []byte("backend:web,team:kenafeh,bleh")
	tags, host := parseTagsStream(rawTags, false)

	assert.Len(t, host, 0)
	assert.Len(t, tags, 3)
}

func BenchmarkParseTagsSplit2(b *testing.B) {
	rawTags := []byte("backend:web,team:kenafeh,bleh")
	for n := 0; n < b.N; n++ {
		_, _ = parseTagsStream(rawTags, true)
	}
}
