package naturalsort

import (
	"bytes"
	"reflect"
	"testing"
)

var assorted = []string{
	"010",
	"00",
	"1",
	"0",
	"01",
	"11",
	"2",
	"A1",
	"A",
	"A11",
	"A2",
	"A11X1",
	"A11X",
	"A11X11",
	"A11X2",
	"Chapter 11 - K",
	"Chapter 2 - B",
	"Chapter 1 - A",
	"v0.1.0",
	"v0.11.0",
	"v0.2.0",
	"v0.2.1",
	"v0.2.11",
	"v0.2.2",
}

var sorted = []string{
	"0",
	"00",
	"01",
	"1",
	"2",
	"010",
	"11",
	"A",
	"A1",
	"A2",
	"A11",
	"A11X",
	"A11X1",
	"A11X2",
	"A11X11",
	"Chapter 1 - A",
	"Chapter 2 - B",
	"Chapter 11 - K",
	"v0.1.0",
	"v0.2.0",
	"v0.2.1",
	"v0.2.2",
	"v0.2.11",
	"v0.11.0",
}

func reverseStringSlice(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func TestReadChunk(t *testing.T) {
	input := "x1yy22,.!abc"
	want := []string{"x", "1", "yy", "22", ",.!abc"}

	var got []string
	reader := bytes.NewReader([]byte(input))
	for {
		chunk := readChunk(reader)
		if chunk == "" {
			break
		}
		got = append(got, chunk)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("invalid chunks; want: %v, got: %v", want, got)
	}
}

func TestCompareNumericChunks(t *testing.T) {
	if compareNumericChunks("1234", "123") != 1 {
		t.Errorf("1234 not greater than 123")
	}
	if compareNumericChunks("123", "1234") != -1 {
		t.Errorf("123 not less than 1234")
	}
	if compareNumericChunks("1234", "1234") != 0 {
		t.Errorf("1234 not equal to 1234")
	}
	// 2**64/2 - 1 < 2**64/2
	if compareNumericChunks("9223372036854775807", "9223372036854775808") != -1 {
		t.Errorf("overflow: numeric comparison fails for big numbers")
	}
	// 2**64 - 1 < 2**64
	if compareNumericChunks("18446744073709551615", "18446744073709551616") != -1 {
		t.Errorf("overflow: numeric comparison fails for big numbers")
	}
}

func TestSort(t *testing.T) {
	input := assorted[:]
	want := sorted

	Sort(input)

	if !reflect.DeepEqual(want, input) {
		t.Errorf("wrong order;\nwant: %v,\n got: %v", want, input)
	}
}

func TestSortReversed(t *testing.T) {
	input := assorted[:]
	want := sorted[:]
	reverseStringSlice(want)

	SortReversed(input)

	if !reflect.DeepEqual(want, input) {
		t.Errorf("wrong order;\nwant: %v,\n got: %v", want, input)
	}
}
