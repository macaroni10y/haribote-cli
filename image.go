package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// PlaceholderConfig holds configuration for generating placeholder images
type PlaceholderConfig struct {
	Width     int
	Height    int
	BgColor   color.RGBA
	TextColor color.RGBA
	Filename  string
}

// GeneratePlaceholderImage creates a placeholder image with the given configuration
func GeneratePlaceholderImage(config PlaceholderConfig) error {
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: config.BgColor}, image.Point{}, draw.Src)

	text := fmt.Sprintf("%d x %d", config.Width, config.Height)

	point := fixed.Point26_6{
		X: fixed.Int26_6((config.Width - len(text)*7) / 2 * 64),
		Y: fixed.Int26_6((config.Height / 2) * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(config.TextColor),
		Face: basicfont.Face7x13, // Using a basic font
		Dot:  point,
	}
	d.DrawString(text)

	outFile, err := os.Create(config.Filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		return fmt.Errorf("failed to encode PNG: %v", err)
	}

	return nil
}
