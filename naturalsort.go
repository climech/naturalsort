// Copyright 2021 climech. All rights reserved. Use of this source code is
// governed by an MIT-style license that can be found in the LICENSE file.

// Package naturalsort implements natural sorting.
package naturalsort

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"
)

func firstRuneIsDigit(scanner io.RuneScanner) (bool, error) {
	r, _, err := scanner.ReadRune()
	if err != nil {
		return false, err
	}
	if err := scanner.UnreadRune(); err != nil {
		return false, err
	}
	return unicode.IsDigit(r), nil
}

// readChunk reads the next numeric or non-numeric chunk from scanner. It
// returns an empty string if there are no more chunks left.
func readChunk(scanner io.RuneScanner) string {
	var chunk strings.Builder

	handleNonEOFError := func(err error) {
		panic(fmt.Sprintf("non-EOF error (invalid UTF-8?): %v", err))
	}

	numeric, err := firstRuneIsDigit(scanner)
	if err != nil {
		if err == io.EOF {
			return ""
		}
		handleNonEOFError(err)
	}

	for {
		r, _, err := scanner.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			handleNonEOFError(err)
		}
		if unicode.IsDigit(r) != numeric {
			if err := scanner.UnreadRune(); err != nil {
				handleNonEOFError(err)
			}
			break
		}
		chunk.WriteRune(r)
	}
	return chunk.String()
}

func chunkIsNumeric(chunk string) bool {
	if len(chunk) != 0 {
		return unicode.IsDigit([]rune(chunk)[0])
	}
	return false
}

func compareNumericChunks(a, b string) int {
	// Pad the smaller number with zeroes and compare the strings
	// lexicographically (slightly more performant than using big.Int).
	la, lb := len(a), len(b)
	if la > lb {
		paddedB := strings.Repeat("0", la-lb) + b
		return strings.Compare(a, paddedB)
	}
	paddedA := strings.Repeat("0", lb-la) + a
	return strings.Compare(paddedA, b)
}

// Compare returns true if a < b according to natural order.
func Compare(a, b string) bool {
	ra := bytes.NewReader([]byte(a))
	rb := bytes.NewReader([]byte(b))

	for {
		chunkA := readChunk(ra)
		chunkB := readChunk(rb)
		if chunkIsNumeric(chunkA) && chunkIsNumeric(chunkB) {
			if cmp := compareNumericChunks(chunkA, chunkB); cmp != 0 {
				return cmp == -1
			} else if cmp := strings.Compare(chunkA, chunkB); cmp != 0 {
				// Zero-padded numbers with unequal number of zeroes -- shortest string
				// takes precedence (same as `sort -V`).
				return cmp == -1
			}
		} else {
			if cmp := strings.Compare(chunkA, chunkB); cmp != 0 {
				return cmp == -1
			}
			if chunkA == "" && chunkB == "" {
				break
			}
		}
	}

	return false
}

// Sort sorts a slice of strings in-place in natural ascending order.
func Sort(slice []string) {
	sort.SliceStable(slice, func(i, j int) bool {
		return Compare(slice[i], slice[j])
	})
}

// Sort sorts a slice of strings in-place in natural descending order.
func SortReversed(slice []string) {
	sort.SliceStable(slice, func(i, j int) bool {
		return Compare(slice[j], slice[i])
	})
}
