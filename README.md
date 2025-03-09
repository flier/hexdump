# hexdump [![Build Status](https://github.com/flier/hexdump/workflows/ci/badge.svg)](https://github.com/flier/hexdump/actions) [![Document](https://pkg.go.dev/badge/github.com/flier/hexdump)](https://pkg.go.dev/github.com/flier/hexdump) [![License](https://img.shields.io/github/license/flier/hexdump)](/LICENSE) [![Release](https://img.shields.io/github/release/flier/hexdump.svg)](https://github.com/flier/hexdump/releases/latest)

Dump binary content as an ASCII table in Golang.

## Install

```sh
go get github.com/flier/hexdump
```

## Examples

Dump string, slices, stream, any value or pointer to value.

### String

```go
hexdump.String("Hello, World!")
// Output:
// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21           |Hello, World!   |
```

### Slices

```go
hexdump.Slices([]uint16{0x0123, 0x4567, 0x89AB, 0xCDEF}, hexdump.TwoBytesHex, hexdump.LittleEndian)
// Output:
// 00000000     0123    4567    89ab    cdef                                  |#.gE....        |
```

### Stream

```go
b := fnv.New128().Sum([]byte("Hello, World!"))

hexdump.Stream(bytes.NewBuffer(b))
// Output:
// 00000000  48 65 6c 6c 6f 2c 20 57  6f 72 6c 64 21 6c 62 27  |Hello, World!lb'|
// 00000010  2e 07 bb 01 42 62 b8 21  75 62 95 c5 8d           |....Bb.!ub...   |
```

### Value

```go
hexdump.Value(complex128(3.14))
// Output:
// 00000000  1f 85 eb 51 b8 1e 09 40  00 00 00 00 00 00 00 00  |...Q...@........|
```

### Structure

```go
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

hexdump.Deref(udp)
// Output:
// 00000000  12 34 56 78 9a bc de f0                           |.4Vx....        |
```

## License

[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](LICENSE) for more details.