package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw" // Standard library draw
	"image/png"
	"log"
	"os"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	xdraw "golang.org/x/image/draw" // For scaling
)

func main() {
	width := flag.Int("width", 0, "The width of the image")
	height := flag.Int("height", 0, "The height of the image")
	bgColorStr := flag.String("bgColor", "grey", "Background color (e.g., grey, #FF0000)")
	textColorStr := flag.String("textColor", "white", "Text color (e.g., white, #00FF00)")
	filename := flag.String("filename", "placeholder.png", "The output filename")
	// Removed fontPath flag

	flag.Parse()

	if *width <= 0 || *height <= 0 {
		fmt.Println("Width and height must be greater than 0")
		flag.Usage()
		os.Exit(1)
	}

	bgColor, err := parseHexColor(*bgColorStr)
	if err != nil {
		bgColor = parseNamedColor(*bgColorStr)
	}

	textColor, err := parseHexColor(*textColorStr)
	if err != nil {
		textColor = parseNamedColor(*textColorStr)
	}

	img := image.NewRGBA(image.Rect(0, 0, *width, *height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	text := fmt.Sprintf("%d x %d", *width, *height)
	
	currentFace := basicfont.Face7x13
	fontMetrics := currentFace.Metrics()

	// 1. Measure text at native size using font.BoundString
	textBounds, _ := font.BoundString(currentFace, text)
	nativeTextRenderWidth := (textBounds.Max.X - textBounds.Min.X).Ceil()
	nativeTextRenderHeight := (textBounds.Max.Y - textBounds.Min.Y).Ceil()

	if nativeTextRenderWidth <= 0 {
		nativeTextRenderWidth = 1 // Avoid zero width
	}
	if nativeTextRenderHeight <= 0 {
		nativeTextRenderHeight = fontMetrics.Height.Ceil()
		if nativeTextRenderHeight <= 0 {
			nativeTextRenderHeight = 13 // Fallback height
		}
	}

	// 2. Create a temporary image for the text
	tempImg := image.NewRGBA(image.Rect(0, 0, nativeTextRenderWidth, nativeTextRenderHeight))
	// Fill with transparent to ensure clean scaling, especially if font has alpha
	draw.Draw(tempImg, tempImg.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// 3. Draw text onto the temporary image
	// Dot for tempDrawer offsets textBounds.Min to (0,0) in tempImg
	tempDrawerDotX := -textBounds.Min.X
	tempDrawerDotY := -textBounds.Min.Y
	
	tempDrawer := &font.Drawer{
		Dst:  tempImg,
		Src:  image.NewUniform(textColor),
		Face: currentFace,
		Dot:  fixed.Point26_6{X: tempDrawerDotX, Y: tempDrawerDotY},
	}
	tempDrawer.DrawString(text)

	// 4. Calculate scale factor and target dimensions
	// Aim for text to occupy ~50% of image width, but not smaller than its native size.
	targetScaledTextWidth := float64(*width) * 0.5 // Changed from 0.7 to 0.5
	if float64(nativeTextRenderWidth) > targetScaledTextWidth {
		targetScaledTextWidth = float64(nativeTextRenderWidth)
	}
	if targetScaledTextWidth < 10 { // Ensure a minimum sensible width
		targetScaledTextWidth = 10
	}

	scaleFactor := 1.0
	if nativeTextRenderWidth > 0 {
		scaleFactor = targetScaledTextWidth / float64(nativeTextRenderWidth)
	}
	if scaleFactor <= 0 {
		scaleFactor = 1.0 // Ensure positive scale
	}

	finalScaledWidth := int(float64(nativeTextRenderWidth) * scaleFactor)
	finalScaledHeight := int(float64(nativeTextRenderHeight) * scaleFactor)

	if finalScaledWidth <= 0 { finalScaledWidth = 1 }
	if finalScaledHeight <= 0 { finalScaledHeight = 1 }

	// 5. Define the destination rectangle on the main image for centering
	dstX := (*width - finalScaledWidth) / 2
	dstY := (*height - finalScaledHeight) / 2
	dstRect := image.Rect(dstX, dstY, dstX+finalScaledWidth, dstY+finalScaledHeight)

	// 6. Scale and draw the temporary image onto the main image
	// ApproxBiLinear for smoother (but potentially blurrier) scaling
	xdraw.ApproxBiLinear.Scale(img, dstRect, tempImg, tempImg.Bounds(), xdraw.Over, nil)
	// For sharper, blockier scaling, use:
	// xdraw.NearestNeighbor.Scale(img, dstRect, tempImg, tempImg.Bounds(), xdraw.Over, nil)

	outFile, err := os.Create(*filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}

	fmt.Printf("Placeholder image '%s' created successfully.\n", *filename)
}

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
