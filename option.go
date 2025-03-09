package hexdump

import (
	"encoding/binary"
	"io"
	"os"
)

// Option can be used to customize the behavior of the [Dumper].
type Option func(*Dumper)

var (
	Stdout = Output(os.Stdout) // Output to stdout.
	Stderr = Output(os.Stderr) // Output to stderr.

	AutoColor   = Color(ColorAuto)   // Auto color mode.
	AlwaysColor = Color(ColorAlways) // Always color mode.
	NeverColor  = Color(ColorNever)  // Never color mode.

	Canonical     = Style(StyleCanonical)     // Canonical hex+ASCII display.
	OneByteChar   = Style(StyleOneByteChar)   // One-byte character display.
	OneByteHex    = Style(StyleOneByteHex)    // One-byte hex display.
	OneByteOctal  = Style(StyleOneByteOctal)  // One-byte octal display.
	TwoBytesDec   = Style(StyleTwoBytesDec)   // Two-byte decimal display.
	TwoBytesHex   = Style(StyleTwoBytesHex)   // Two-byte hexadecimal display
	TwoBytesOctal = Style(StyleTwoBytesOctal) // Wwo-byte octal display

	LittleEndian = ByteOrder(binary.LittleEndian) // Little-endian byte order.
	BigEndian    = ByteOrder(binary.BigEndian)    // Big-endian byte order.
	NativeEndian = ByteOrder(binary.NativeEndian) // Native-endian byte order.
)

// The output stream, the default is [os.Stdout].
func Output(w io.Writer) Option { return func(d *Dumper) { d.Output = w } }

// The number of bytes per line, the default is 16.
func LineWidth(n int) Option { return func(d *Dumper) { d.LineWidth = n } }

// The color mode of the line, the default is [ColorAuto].
func Color(c ColorMode) Option { return func(d *Dumper) { d.Color = c } }

// The color theme of the line, the default is [DefaultTheme].
func Theme(t *ColorTheme) Option { return func(d *Dumper) { d.Theme = t } }

// The display style of the line, the default is [StyleCanonical].
func Style(s DisplayStyle) Option { return func(d *Dumper) { d.Style = s } }

// The byte order used to read the data group, the default is [binary.NativeEndian].
func ByteOrder(b binary.ByteOrder) Option { return func(d *Dumper) { d.ByteOrder = b } }

// The start offset of the binary content.
func Start(off int64) Option { return func(d *Dumper) { d.Start = off } }

// Skip offset bytes from the beginning of the input.
func Skip(n int64) Option { return func(d *Dumper) { d.Skip = n } }

// Interpret only length bytes of input.
func Length(n int64) Option { return func(d *Dumper) { d.Length = n } }

// Extract the range of input from start to end.
func Range(start, end int64) Option {
	if start > end {
		start, end = end, start
	}

	return func(d *Dumper) {
		d.Skip = start
		d.Length = end - start
	}
}
