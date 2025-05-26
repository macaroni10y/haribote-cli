package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	width := flag.Int("width", 0, "The width of the image")
	height := flag.Int("height", 0, "The height of the image")
	bgColorStr := flag.String("bgColor", "grey", "Background color (e.g., grey, #FF0000)")
	textColorStr := flag.String("textColor", "white", "Text color (e.g., white, #00FF00)")
	filename := flag.String("filename", "placeholder.png", "The output filename")

	flag.Parse()

	if *width <= 0 || *height <= 0 {
		fmt.Println("Width and height must be greater than 0")
		flag.Usage()
		os.Exit(1)
	}

	bgColor := parseColor(*bgColorStr)
	textColor := parseColor(*textColorStr)

	config := PlaceholderConfig{
		Width:     *width,
		Height:    *height,
		BgColor:   bgColor,
		TextColor: textColor,
		Filename:  *filename,
	}

	err := GeneratePlaceholderImage(config)
	if err != nil {
		log.Fatalf("Failed to generate placeholder image: %v", err)
	}

	fmt.Printf("Placeholder image '%s' created successfully.\n", *filename)
}
