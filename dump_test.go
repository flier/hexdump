package hexdump_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/hexdump"
)

func ExampleString() {
	_ = hexdump.String("Hello, World!")
	// Output:
	// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21           |Hello, World!   |
}

func ExampleSlices() {
	_ = hexdump.Slices([]uint16{0x0123, 0x4567, 0x89AB, 0xCDEF}, hexdump.TwoBytesHex, hexdump.LittleEndian)
	// Output:
	// 00000000     0123    4567    89ab    cdef                                  |#.gE....        |
}

func ExampleBytes() {
	_ = hexdump.Bytes(fnv.New128().Sum([]byte("Hello, World!")))
	// Output:
	// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21 6c 62 27  |Hello, World!lb'|
	// 00000010  2e 07 bb 01 42 62 b8 21  75 62 95 c5 8d           |....Bb.!ub...   |
}

func ExampleStream() {
	b := fnv.New128().Sum([]byte("Hello, World!"))

	_ = hexdump.Stream(bytes.NewBuffer(b))
	// Output:
	// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21 6c 62 27  |Hello, World!lb'|
	// 00000010  2e 07 bb 01 42 62 b8 21  75 62 95 c5 8d           |....Bb.!ub...   |
}

func ExampleValue() {
	_ = hexdump.Value(complex128(3.14))
	// Output:
	// 00000000  1f 85 eb 51 b8 1e 09 40  00 00 00 00 00 00 00 00  |...Q...@........|
}

func ExampleDeref() {
	type UDP struct {
		SrcPort uint16
		DstPort uint16
		Length  uint16
		Chksum  uint16
	}

	htons := func(v uint16) uint16 {
		var b [2]byte

		binary.NativeEndian.PutUint16(b[:], v)

		return binary.BigEndian.Uint16(b[:])
	}

	udp := &UDP{
		SrcPort: htons(0x1234),
		DstPort: htons(0x5678),
		Length:  htons(0x9abc),
		Chksum:  htons(0xdef0),
	}

	_ = hexdump.Deref(udp)
	// Output:
	// 00000000  12 34 56 78 9a bc de f0                           |.4Vx....        |
}

func ExampleAsDeref() {
	var b [8]byte

	binary.BigEndian.PutUint64(b[:], 0x123456789abcdef)

	_ = hexdump.AsDeref[uint64](&b[0])
	// Output:
	// 00000000  01 23 45 67 89 ab cd ef                           |.#Eg....        |
}

func ExampleOutput() {
	var b strings.Builder

	_ = hexdump.String("Hello, World!", hexdump.Output(&b))

	fmt.Print(b.String())
	// Output:
	// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21           |Hello, World!   |
}

func ExampleLineWidth() {
	_ = hexdump.String("Hello, World!", hexdump.LineWidth(8))
	// Output:
	// 00000000  48 65 6c 6c 6f 2c 20 57  |Hello, W|
	// 00000008  6f 72 6c 64 21           |orld!   |
}

func ExampleCanonical() {
	_ = hexdump.String("Hello, World!", hexdump.Canonical)
	// Output:
	// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21           |Hello, World!   |
}

func ExampleOneByteHex() {
	_ = hexdump.String("Hello, World!", hexdump.OneByteHex)
	// Output:
	// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21           |Hello, World!   |
}

func ExampleOneByteOctal() {
	_ = hexdump.String("Hello, World!", hexdump.OneByteOctal)
	// Output:
	// 00000000  110 145 154 154 157 054 040 127  157 162 154 144 041              |Hello, World!   |
}

func ExampleOneByteChar() {
	_ = hexdump.String("Hello, World!", hexdump.OneByteChar)
	// Output:
	// 00000000    H   e   l   l   o   ,       W    o   r   l   d   !              |Hello, World!   |
}

func ExampleTwoBytesDec() {
	_ = hexdump.Slices([]uint16{123, 456, 789}, hexdump.TwoBytesDec)
	// Output:
	// 00000000    00123   00456   00789                                          |{.....          |
}

func ExampleTwoBytesHex() {
	_ = hexdump.Slices([]uint16{0x0123, 0x4567, 0x89ab, 0xcdef}, hexdump.TwoBytesHex, hexdump.LittleEndian)
	// Output:
	// 00000000     0123    4567    89ab    cdef                                  |#.gE....        |
}

func ExampleTwoBytesOctal() {
	_ = hexdump.Slices([]uint16{0o123, 0o456, 0o777}, hexdump.TwoBytesOctal)
	// Output:
	// 00000000   000123  000456  000777                                          |S.....          |
}

func ExampleLittleEndian() {
	_ = hexdump.Bytes([]byte{1, 2, 3, 4, 5, 6, 7, 8}, hexdump.TwoBytesHex, hexdump.LittleEndian)
	// Output:
	// 00000000     0201    0403    0605    0807                                  |........        |
}

func ExampleBigEndian() {
	_ = hexdump.Bytes([]byte{1, 2, 3, 4, 5, 6, 7, 8}, hexdump.TwoBytesHex, hexdump.BigEndian)
	// Output:
	// 00000000     0102    0304    0506    0708                                  |........        |
}

func ExampleNativeEndian() {
	_ = hexdump.Slices([]uint16{0x0123, 0x4567, 0x89ab, 0xcdef}, hexdump.TwoBytesHex, hexdump.NativeEndian)
	// Output:
	// 00000000     0123    4567    89ab    cdef                                  |#.gE....        |
}

func ExampleStart() {
	_ = hexdump.String("Hello, World!", hexdump.Start(0x1004))
	// Output:
	// 00001000              48 65 6c 6c  6f 2c 20 57 6f 72 6c 64  |    Hello, World|
	// 00001010  21                                                |!               |
}

func ExampleSkip() {
	_ = hexdump.String("Hello, World!", hexdump.Skip(7))
	// Output:
	// 00000000                       57  6f 72 6c 64 21           |       World!   |
}

func ExampleSkip_length() {
	_ = hexdump.String("Hello, World!", hexdump.Skip(7), hexdump.Length(5))
	// Output:
	// 00000000                       57  6f 72 6c 64              |       World    |
}

func ExampleRange() {
	_ = hexdump.String("Hello, World!", hexdump.Range(7, 12))
	// Output:
	// 00000000                       57  6f 72 6c 64              |       World    |
}

func ExampleLength() {
	_ = hexdump.String("Hello, World!", hexdump.Length(5))
	// Output:
	// 00000000  48 65 6c 6c 6f                                    |Hello           |
}

func TestAlwaysColor(t *testing.T) {
	t.Parallel()

	Convey("Given a some string", t, func() {
		s := "Hello, World!"

		Convey("When dump it with the default color theme", func() {
			var b strings.Builder

			_ = hexdump.String(s, hexdump.AlwaysColor, hexdump.Output(&b))

			Convey("Then the output should contains color", func() {
				theme := hexdump.DefaultTheme

				So(b.String(), ShouldEqual,
					theme.Offset.Sprint("00000000")+" "+
						theme.Content.Sprint(" 48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21         ")+"  |"+
						theme.Chars.Sprint("Hello, World!   ")+"|\n")
			})
		})
	})
}
