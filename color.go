package main

import (
	"fmt"
	"image/color"
	"strconv"
)

// parseHexColor parses a hex color string (e.g., "#RRGGBB" or "RRGGBB")
func parseHexColor(s string) (color.RGBA, error) {
	if s[0] == '#' {
		s = s[1:]
	}
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color string: %s", s)
	}
	r, err := strconv.ParseInt(s[0:2], 16, 0)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseInt(s[2:4], 16, 0)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseInt(s[4:6], 16, 0)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}

// parseNamedColor provides a small map of named colors.
// For a more comprehensive solution, a larger library or map would be needed.
func parseNamedColor(s string) color.RGBA {
	switch s {
	case "black":
		return color.RGBA{0, 0, 0, 255}
	case "white":
		return color.RGBA{255, 255, 255, 255}
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "green":
		return color.RGBA{0, 255, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "grey", "gray":
		return color.RGBA{128, 128, 128, 255}
	default:
		fmt.Printf("Warning: Color '%s' not recognized, defaulting to grey.\n", s)
		return color.RGBA{128, 128, 128, 255}
	}
}

// parseColor attempts to parse a color string as hex first, then falls back to named colors
func parseColor(colorStr string) color.RGBA {
	if color, err := parseHexColor(colorStr); err == nil {
		return color
	}
	return parseNamedColor(colorStr)
}
