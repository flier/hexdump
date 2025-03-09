package hexdump

import (
	"encoding/binary"
	"fmt"
	"iter"
	"slices"
	"unicode"
)

// The display style of binary content.
type DisplayStyle int

const (
	StyleCanonical     DisplayStyle = iota // Canonical hex+ASCII display.
	StyleOneByteChar                       // One-byte character display.
	StyleOneByteHex                        // One-byte hex display.
	StyleOneByteOctal                      // One-byte octal display.
	StyleTwoBytesDec                       // Two-byte decimal display.
	StyleTwoBytesHex                       // Two-byte hexadecimal display
	StyleTwoBytesOctal                     // Wwo-byte octal display
)

func (s DisplayStyle) formatLine(width, skip int, buf []byte, order binary.ByteOrder) []string {
	return slices.Concat(
		s.padding(max(skip, 0)),
		slices.Collect(s.formatGroups(buf, order)),
		s.padding(max(width-skip-len(buf), 0)))
}

const twoBytes = 2

func (s DisplayStyle) padding(n int) []string {
	switch s {
	case StyleCanonical, StyleOneByteHex:
		return slices.Repeat([]string{"  "}, n)

	case StyleOneByteChar, StyleOneByteOctal:
		return slices.Repeat([]string{"   "}, n)

	case StyleTwoBytesHex, StyleTwoBytesOctal, StyleTwoBytesDec:
		return slices.Repeat([]string{"       "}, n/twoBytes)

	default:
		return []string{}
	}
}

func (s DisplayStyle) formatGroups(b []byte, order binary.ByteOrder) iter.Seq[string] {
	return func(yield func(s string) bool) {
		rest := b

		for len(rest) > 0 {
			var str string

			str, rest = s.formatGroup(rest, order)

			if !yield(str) {
				return
			}
		}
	}
}

func (s DisplayStyle) formatGroup(b []byte, order binary.ByteOrder) (text string, rest []byte) {
	switch s {
	case StyleCanonical, StyleOneByteHex, StyleOneByteOctal:
		return s.formatValue(b[0]), b[1:]

	case StyleOneByteChar:
		if unicode.IsPrint(rune(b[0])) {
			return s.formatValue(b[0]), b[1:]
		}

		return s.padding(1)[0], b[1:]

	case StyleTwoBytesDec, StyleTwoBytesHex, StyleTwoBytesOctal:
		var v uint16

		if len(b) == 1 {
			v, rest = uint16(b[0]), nil
		} else {
			v, rest = order.Uint16(b[:2]), b[2:]
		}

		text = s.formatValue(v)

		return

	default:
		panic(s)
	}
}

func (s DisplayStyle) formatValue(v ...any) string {
	return fmt.Sprintf(formatValues[s], v...)
}

var formatValues = map[DisplayStyle]string{
	StyleCanonical:     "%02x",
	StyleOneByteChar:   "  %c",
	StyleOneByteHex:    "%02x",
	StyleOneByteOctal:  "%03o",
	StyleTwoBytesDec:   "  %05d",
	StyleTwoBytesHex:   "   %04x",
	StyleTwoBytesOctal: " %06o",
}
