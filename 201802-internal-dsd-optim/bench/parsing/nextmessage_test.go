package parsing

import (
	"bufio"
	"bytes"
	"testing"
)

const packet = "daemon:666|g|#sometag1:somevalue1,sometag2:somevalue2\ndaemon:666|g|\n \n"

func BenchmarkNextMessageCurrent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b := []byte(packet)
		var i int

		for {
			split := bytes.SplitAfterN(b, []byte("\n"), 2)
			b = b[len(split[0]):]
			// Remove trailing newline
			if len(split) == 2 {
				i++
			} else {
				break
			}

		}
		if i != 3 {
			panic(i)
		}
	}
}

func BenchmarkNextMessageScanner(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b := []byte(packet)
		var i int

		reader := bytes.NewReader(b)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			_ = scanner.Bytes()
			i++
		}
		if i != 3 {
			panic(i)
		}
	}
}

func BenchmarkNextMessageScanLines(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b := []byte(packet)
		var i int

		for {
			advance, line, err := bufio.ScanLines(b, false)
			if err != nil || len(line) == 0 {
				break
			}
			b = b[advance:]
			i++
		}
		if i != 3 {
			panic(i)
		}
	}
}
