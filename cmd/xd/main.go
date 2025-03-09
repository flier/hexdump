package main

import (
	"flag"
	"io"
	"log/slog"
	"os"

	. "github.com/flier/hexdump" //nolint:revive,stylecheck
)

var (
	canonical    = flag.Bool("C", false, "canonical hex+ASCII display")
	oneByteOctal = flag.Bool("b", false, "one-byte octal")
	oneByteHex   = flag.Bool("X", false, "one-byte hex")
	oneByteChar  = flag.Bool("c", false, "one-byte char")
	twoBytesDec  = flag.Bool("d", false, "two-byte decimal")
	twoBytesOct  = flag.Bool("e", false, "two-byte octal")
	twoBytesHex  = flag.Bool("x", false, "two-byte hex")
	color        = ColorAuto
	noColor      = flag.Bool("no-color", false, "disable color mode")
	length       = flag.Int64("n", 0, "interpret only length bytes of input")
	skip         = flag.Int64("s", 0, "skip first skip bytes of input")
	width        = flag.Int("w", DefaultLineWidth, "output line width")
	verbose      = flag.Bool("v", false, "show verbose messages")
	debug        = flag.Bool("vv", false, "show debug messages")
)

func main() {
	flag.TextVar(&color, "L", color, "color mode")
	flag.Parse()

	initLogger()

	if flag.NArg() == 0 {
		dump("-", os.Stdin)
	} else {
		for _, name := range flag.Args() {
			f, err := os.Open(name)
			if err != nil {
				slog.Warn("open file", "err", err)
			}

			dump(name, f)
		}
	}
}

func initLogger() {
	switch {
	case *debug:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case *verbose:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}

func dump(name string, r io.Reader) {
	opts := []Option{
		Style(displayStyle()),
		Color(colorMode()),
		Length(*length),
		Skip(*skip),
		LineWidth(*width),
	}

	err := Stream(r, opts...)
	if err != nil {
		slog.Error("hexdump stream", "name", name, "err", err)
	}
}

func displayStyle() DisplayStyle {
	switch {
	case *canonical:
		return StyleCanonical
	case *oneByteChar:
		return StyleOneByteChar
	case *oneByteHex:
		return StyleOneByteHex
	case *oneByteOctal:
		return StyleOneByteOctal
	case *twoBytesDec:
		return StyleTwoBytesDec
	case *twoBytesHex:
		return StyleTwoBytesHex
	case *twoBytesOct:
		return StyleTwoBytesOctal
	default:
		return StyleCanonical
	}
}

func colorMode() ColorMode {
	if *noColor {
		return ColorNever
	}

	return color
}
