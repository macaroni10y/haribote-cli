# haribote-cli
[![GitHub License](https://img.shields.io/github/license/macaroni10y/haribote-cli)](https://github.com/macaroni10y/haribote-cli/blob/main/LICENSE)
[![Author](https://img.shields.io/badge/Author-macaroni10y-pink)](https://github.com/macaroni10y)

A simple Go CLI tool for generating placeholder images.

## Overview

haribote-cli is a command-line tool that quickly generates placeholder images with specified dimensions and colors. It's useful when you need temporary image files during web development or prototyping.

## Features

- Generate images of any size
- Customizable background and text colors
- Support for both named colors (red, blue, grey, etc.) and HEX colors (#FF0000)
- PNG format output
- Automatically draws size information on the image

## Installation

### Prerequisites

- Go 1.24.3 or later

### Build

```bash
git clone https://github.com/macaroni10y/haribote-cli
cd haribote-cli
go build -o haribote
```

## Usage

```bash
./haribote -width [WIDTH] -height [HEIGHT] [OPTIONS]
```

### Required Parameters

- `-width`: Image width in pixels
- `-height`: Image height in pixels

### Optional Parameters

- `-bgColor`: Background color (default: "grey")
- `-textColor`: Text color (default: "white")
- `-filename`: Output filename (default: "placeholder.png")

## Examples

### Basic Usage

```bash
# Generate an 800x600 image with grey background
./haribote -width 800 -height 600
```

### Custom Colors

```bash
# Generate an image with blue background and white text
./haribote -width 400 -height 300 -bgColor blue -textColor white
```

### Using HEX Colors

```bash
# Generate an image using HEX colors
./haribote -width 1920 -height 1080 -bgColor "#FF5733" -textColor "#FFFFFF"
```

### Custom Output Filename

```bash
# Output with custom filename
./haribote -width 500 -height 500 -filename custom_placeholder.png
```

## Supported Colors

### Named Colors

- black
- white
- red
- green
- blue
- grey/gray

### HEX Colors

Supports 6-digit HEX color codes in either `#RRGGBB` or `RRGGBB` format.

## Technical Specifications

- **Language**: Go
- **Output Format**: PNG
- **Dependencies**: 
  - `golang.org/x/image` - Image processing and font rendering
