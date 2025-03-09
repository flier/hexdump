package hexdump

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type Formatter struct {
	*bufio.Writer
	*ColorTheme
	DisplayStyle
	binary.ByteOrder
	LineWidth int
}

const groupsSep = 8

func (f *Formatter) FormatLine(off int64, skip int, buf []byte) (err error) {
	return errors.Join(
		f.formatOffset(off),
		f.formatContent(skip, buf),
		f.formatChars(skip, buf),
		f.Flush())
}

func (f *Formatter) formatOffset(off int64) (err error) {
	offset := f.Offset.Sprint(fmt.Sprintf("%08x", off))

	_, err = f.WriteString(offset + " ")

	return
}

func (f *Formatter) formatContent(skip int, buf []byte) (err error) {
	f.Content.SetWriter(f.Writer)
	defer f.Content.UnsetWriter(f.Writer)

	for i, s := range f.formatLine(f.LineWidth, skip, buf, f.ByteOrder) {
		if err = f.WriteByte(' '); err != nil {
			return
		}

		if i > 0 && f.LineWidth > groupsSep && i%groupsSep == 0 {
			if err = f.WriteByte(' '); err != nil {
				return
			}
		}

		if _, err = f.WriteString(s); err != nil {
			return
		}
	}

	return
}

func (f *Formatter) formatChars(skip int, buf []byte) (err error) {
	chars := f.Chars.Sprint(f.charTable(skip, buf))

	_, err = f.WriteString("  |" + chars + "|\n")

	return
}

func (f *Formatter) charTable(skip int, buf []byte) string {
	return spaces(skip) +
		string(slices.Collect(func(yield func(byte) bool) {
			for _, c := range buf {
				if c >= 0x80 || !unicode.IsPrint(rune(c)) {
					c = '.'
				}

				if !yield(c) {
					return
				}
			}
		})) + spaces(f.LineWidth-skip-len(buf))
}

func spaces(n int) string {
	return strings.Repeat(" ", max(n, 0))
}
