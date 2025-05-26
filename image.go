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
	xdraw "golang.org/x/image/draw"
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
	currentFace := basicfont.Face7x13

	textBounds, _ := font.BoundString(currentFace, text)
	nativeTextRenderWidth := (textBounds.Max.X - textBounds.Min.X).Ceil()
	nativeTextRenderHeight := (textBounds.Max.Y - textBounds.Min.Y).Ceil()

	if nativeTextRenderWidth <= 0 {
		nativeTextRenderWidth = 1
	}
	if nativeTextRenderHeight <= 0 {
		fontMetrics := currentFace.Metrics()
		nativeTextRenderHeight = fontMetrics.Height.Ceil()
		if nativeTextRenderHeight <= 0 {
			nativeTextRenderHeight = 13 // Fallback to basicfont height
		}
	}

	tempImg := image.NewRGBA(image.Rect(0, 0, nativeTextRenderWidth, nativeTextRenderHeight))
	xdraw.Draw(tempImg, tempImg.Bounds(), image.Transparent, image.Point{}, xdraw.Src)

	tempDrawerDotX := -textBounds.Min.X
	tempDrawerDotY := -textBounds.Min.Y
	tempDrawer := &font.Drawer{
		Dst:  tempImg,
		Src:  image.NewUniform(config.TextColor),
		Face: currentFace,
		Dot:  fixed.Point26_6{X: tempDrawerDotX, Y: tempDrawerDotY},
	}
	tempDrawer.DrawString(text)

	targetScaledTextWidth := float64(config.Width) * 0.7
	if float64(nativeTextRenderWidth) > targetScaledTextWidth {
		targetScaledTextWidth = float64(nativeTextRenderWidth)
	}
	if targetScaledTextWidth < 10 {
		targetScaledTextWidth = 10
	}

	scaleFactor := 1.0
	if nativeTextRenderWidth > 0 {
		scaleFactor = targetScaledTextWidth / float64(nativeTextRenderWidth)
	}
	if scaleFactor <= 0 {
		scaleFactor = 1.0
	}

	finalScaledWidth := int(float64(nativeTextRenderWidth) * scaleFactor)
	finalScaledHeight := int(float64(nativeTextRenderHeight) * scaleFactor)

	if finalScaledWidth <= 0 {
		finalScaledWidth = 1
	}
	if finalScaledHeight <= 0 {
		finalScaledHeight = 1
	}

	dstX := (config.Width - finalScaledWidth) / 2
	dstY := (config.Height - finalScaledHeight) / 2
	dstRect := image.Rect(dstX, dstY, dstX+finalScaledWidth, dstY+finalScaledHeight)

	xdraw.ApproxBiLinear.Scale(img, dstRect, tempImg, tempImg.Bounds(), xdraw.Over, nil)

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
