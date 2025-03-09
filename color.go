package hexdump

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

//go:generate go tool stringer -type=ColorMode -linecomment

type ColorMode int //nolint:recvcheck

const (
	ColorAuto   ColorMode = iota // auto
	ColorAlways                  // always
	ColorNever                   // never
)

func (m ColorMode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *ColorMode) UnmarshalText(text []byte) (err error) {
	for i := range len(_ColorMode_index) - 1 {
		if strings.EqualFold(string(text), _ColorMode_name[_ColorMode_index[i]:_ColorMode_index[i+1]]) {
			*m = ColorMode(i)

			return nil
		}
	}

	return fmt.Errorf("color mode %q, %w", text, os.ErrInvalid)
}

func initColor(w io.Writer, m ColorMode) {
	switch m {
	case ColorAuto:
		color.NoColor = noColorIsSet() || termIsDumb() || !isTty(w)

	case ColorAlways:
		color.NoColor = false

	case ColorNever:
		color.NoColor = true
	}
}

func noColorIsSet() bool { return os.Getenv("NO_COLOR") != "" }
func termIsDumb() bool   { return os.Getenv("TERM") == "dumb" }

type file interface {
	Fd() uintptr
}

func isTty(w io.Writer) bool {
	if f, ok := w.(file); ok {
		fd := f.Fd()

		return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
	}

	return false
}

type ColorTheme struct {
	Offset  *color.Color
	Content *color.Color
	Chars   *color.Color
}

var DefaultTheme = ColorTheme{
	Offset:  color.New(color.Faint),
	Content: color.New(color.Reset),
	Chars:   color.New(color.Italic),
}
