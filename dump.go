// Display binary contents in hexadecimal, decimal, octal, or ascii.
package hexdump

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"iter"
	"os"
	"sync"
	"unsafe"
)

const DefaultLineWidth = 16

// Bytes converts a byte slice into a readable ASCII table using the provided options.
// It writes the output to the specified [io.Writer] (default is [os.Stdout]).
//
// The function accepts a byte slice 'b' and a variable number of options 'x'.
// The options can be used to customize the behavior of the dumper.
//
// The function returns an error if any error occurs during the dumping process.
func Bytes(b []byte, x ...Option) (err error) {
	d := New(x...)

	if d.Skip > 0 {
		b = b[d.Skip:]
		d.off = d.Skip
	}

	if d.Length > 0 {
		b = b[:d.Length]
	}

	if _, err = d.Write(b); err != nil {
		return
	}

	return d.Flush()
}

// Slices converts a slice of T into a readable ASCII table using the provided options.
func Slices[T any](s []T, x ...Option) error {
	var v T

	b := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(s))), int(unsafe.Sizeof(v))*len(s))

	return Bytes(b, x...)
}

// String converts a string into a readable ASCII table using the provided options.
func String(s string, x ...Option) error {
	return Bytes([]byte(s), x...)
}

// Value converts the contents of the value 'v' into a readable ASCII table using the options provided.
func Value[T any](v T, x ...Option) (err error) {
	b := unsafe.Slice((*byte)(unsafe.Pointer(&v)), int(unsafe.Sizeof(v)))

	return Bytes(b, x...)
}

// Deref converts the contents of the value pointed to by 'p' into a readable ASCII table using the options provided.
func Deref[T any](p *T, x ...Option) (err error) {
	b := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))

	return Bytes(b, x...)
}

// AsDeref converts the contents of the value pointed to by 'p' into a readable ASCII table using the options provided.
func AsDeref[T, S any](p *S, x ...Option) (err error) {
	return Deref((*T)(unsafe.Pointer(p)), x...)
}

// Stream reads from the provided [io.Reader] and converts the binary content into a readable ASCII table.
// The function writes the output to the specified [io.Writer] (default is [os.Stdout]).
//
// The function accepts a byte stream 'r' and a variable number of options 'x'.
// The options can be used to customize the behavior of the dumper.
//
// The function returns an error if any error occurs during the dumping process.
func Stream(r io.Reader, x ...Option) (err error) {
	d := New(x...)

	if d.Skip > 0 {
		if d.off, err = io.CopyN(io.Discard, r, d.Skip); err != nil {
			return
		}
	}

	if d.Length > 0 {
		r = &io.LimitedReader{R: r, N: d.Length}
	}

	if _, err = io.Copy(d, r); err != nil {
		return
	}

	return d.Flush()
}

// Seq reads value from [iter.Seq] and converts it as binary content into a readable ASCII table.
// The function writes the output to the specified [io.Writer] (default is [os.Stdout]).
//
// The function accepts an iterator 'i' and a variable number of options 'x'.
// The options can be used to customize the behavior of the dumper.
//
// The function returns an error if any error occurs during the dumping process.
func Seq[T any](i iter.Seq[T], x ...Option) (err error) {
	d := New(x...)

	for v := range i {
		b := unsafe.Slice((*byte)(unsafe.Pointer(&v)), int(unsafe.Sizeof(v)))

		if _, err = d.Write(b); err != nil {
			return
		}
	}

	return d.Flush()
}

// Dumper converts the binary content into a readable ASCII table.
type Dumper struct {
	b    bytes.Buffer
	f    *Formatter
	once sync.Once
	off  int64

	// The output stream, the default is [os.Stdout].
	Output io.Writer

	// The number of bytes per line, the default is 16.
	LineWidth int

	// The color mode of the line, the default is [ColorAuto].
	Color ColorMode

	// The color theme of the line, the default is [DefaultTheme].
	Theme *ColorTheme

	// The display style of the line, the default is [StyleCanonical].
	Style DisplayStyle

	// The byte order used to read the data group, the default is [binary.NativeEndian].
	ByteOrder binary.ByteOrder

	// The start offset of the binary content.
	Start int64

	// Skip offset bytes from the beginning of the input.
	Skip int64

	// Interpret only length bytes of input.
	Length int64
}

// New returns a new [Dumper] with the provided options.
func New(x ...Option) (d *Dumper) {
	d = new(Dumper)

	for _, opt := range x {
		opt(d)
	}

	return
}

// Flush dump any buffered data to the underlying [io.Writer].
func (d *Dumper) Flush() (err error) {
	d.init()

	if err = d.flushLines(true); err != nil {
		return err
	}

	return
}

// Write writes the contents of p into the buffer.
//
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining why the write is short.
func (d *Dumper) Write(p []byte) (nn int, err error) {
	if nn, err = d.b.Write(p); err != nil {
		return
	}

	if err = d.flushLines(false); err != nil {
		return 0, err
	}

	return
}

// WriteByte writes a single byte.
func (d *Dumper) WriteByte(c byte) (err error) {
	if err = d.b.WriteByte(c); err != nil {
		return
	}

	if err = d.flushLines(false); err != nil {
		return
	}

	return
}

// WriteRune writes a single Unicode code point, returning the number of bytes written and any error.
func (d *Dumper) WriteRune(r rune) (size int, err error) {
	if size, err = d.b.WriteRune(r); err != nil {
		return
	}

	if err = d.flushLines(false); err != nil {
		return 0, err
	}

	return
}

// WriteString writes a string.
//
// It returns the number of bytes written.
// If the count is less than len(s), it also returns an error explaining why the write is short.
func (d *Dumper) WriteString(s string) (count int, err error) {
	if count, err = d.b.WriteString(s); err != nil {
		return
	}

	if err = d.flushLines(false); err != nil {
		return 0, err
	}

	return
}

func (d *Dumper) flushLines(all bool) (err error) {
	d.once.Do(d.init)

	width := int64(d.LineWidth)
	off := d.Start + d.off
	skip := off % width

	for int(skip)+d.b.Len() >= d.LineWidth {
		if err = d.flushLine(false); err != nil {
			return
		}
	}

	if all && d.b.Len() > 0 {
		if err = d.flushLine(true); err != nil {
			return
		}
	}

	return
}

func (d *Dumper) flushLine(all bool) (err error) {
	width := int64(d.LineWidth)
	off := d.Start + d.off
	start := off / width * width
	skip := off % width

	length := width - skip
	if all {
		length = int64(d.b.Len())
	}

	b := make([]byte, length)

	var n int

	if n, err = io.ReadFull(&d.b, b); err != nil {
		return
	}

	if err = d.f.FormatLine(start, int(skip), b); err != nil {
		return
	}

	d.off += int64(n)

	return
}

func (d *Dumper) init() {
	if d.Output == nil {
		d.Output = os.Stdout
	}

	initColor(d.Output, d.Color)

	if d.Theme == nil {
		d.Theme = &DefaultTheme
	}

	if d.ByteOrder == nil {
		d.ByteOrder = binary.NativeEndian
	}

	if d.LineWidth == 0 {
		d.LineWidth = DefaultLineWidth
	}

	if d.f == nil {
		d.f = &Formatter{bufio.NewWriter(d.Output), d.Theme, d.Style, d.ByteOrder, d.LineWidth}
	}
}
