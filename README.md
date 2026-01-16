# imgutils-watermark

[![Go Reference](https://pkg.go.dev/badge/github.com/imgutils-org/imgutils-watermark.svg)](https://pkg.go.dev/github.com/imgutils-org/imgutils-watermark)
[![Go Report Card](https://goreportcard.com/badge/github.com/imgutils-org/imgutils-watermark)](https://goreportcard.com/report/github.com/imgutils-org/imgutils-watermark)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library for adding watermarks to images. Part of the [imgutils](https://github.com/imgutils-org) collection.

## Features

- Position-based watermark placement (center, corners)
- Configurable opacity
- Tiled watermark patterns
- Padding and margin controls
- Support for PNG watermarks with transparency

## Installation

```bash
go get github.com/imgutils-org/imgutils-watermark
```

## Quick Start

```go
package main

import (
    "image"
    "log"
    "os"

    "github.com/imgutils-org/imgutils-watermark"
)

func main() {
    // Load source image
    srcFile, _ := os.Open("photo.jpg")
    src, _, _ := image.Decode(srcFile)
    srcFile.Close()

    // Load watermark (usually a PNG with transparency)
    wmFile, _ := os.Open("logo.png")
    wm, _, _ := image.Decode(wmFile)
    wmFile.Close()

    // Apply watermark in bottom-right corner
    result := watermark.Apply(src, wm, watermark.DefaultOptions())

    // Save result
    out, _ := os.Create("watermarked.jpg")
    watermark.SaveJPEG(result, out, 90)
    out.Close()
}
```

## Usage Examples

### Basic Watermarking

```go
// Default: bottom-right corner, 50% opacity, 10px padding
result := watermark.Apply(src, wm, watermark.DefaultOptions())
```

### Position Options

```go
// Center
opts := watermark.Options{Position: watermark.Center, Opacity: 0.5}
result := watermark.Apply(src, wm, opts)

// Top-left corner
opts := watermark.Options{Position: watermark.TopLeft, Opacity: 0.5}
result := watermark.Apply(src, wm, opts)

// Top-right corner
opts := watermark.Options{Position: watermark.TopRight, Opacity: 0.5}
result := watermark.Apply(src, wm, opts)

// Bottom-left corner
opts := watermark.Options{Position: watermark.BottomLeft, Opacity: 0.5}
result := watermark.Apply(src, wm, opts)

// Bottom-right corner (default)
opts := watermark.Options{Position: watermark.BottomRight, Opacity: 0.5}
result := watermark.Apply(src, wm, opts)
```

### Opacity Control

```go
// Subtle watermark (30% opacity)
opts := watermark.Options{
    Position: watermark.BottomRight,
    Opacity:  0.3,
}

// Bold watermark (80% opacity)
opts := watermark.Options{
    Position: watermark.BottomRight,
    Opacity:  0.8,
}
```

### Custom Padding

```go
opts := watermark.Options{
    Position: watermark.BottomRight,
    Opacity:  0.5,
    PaddingX: 20, // 20px from right edge
    PaddingY: 20, // 20px from bottom edge
}
```

### Tiled Watermark

```go
// Cover entire image with repeated watermark
result := watermark.Tile(src, wm, 0.3, 50)
// 0.3 = 30% opacity
// 50 = 50px spacing between tiles
```

### From File Paths

```go
result, err := watermark.ApplyFromFiles("photo.jpg", "logo.png", watermark.Options{
    Position: watermark.BottomRight,
    Opacity:  0.5,
    PaddingX: 15,
    PaddingY: 15,
})
if err != nil {
    log.Fatal(err)
}
```

## API Reference

### Types

#### Position

```go
type Position int

const (
    Center      Position = iota // Center of image
    TopLeft                     // Top-left corner
    TopRight                    // Top-right corner
    BottomLeft                  // Bottom-left corner
    BottomRight                 // Bottom-right corner
)
```

#### Options

```go
type Options struct {
    Position Position // Where to place watermark
    Opacity  float64  // 0.0 to 1.0
    PaddingX int      // Horizontal padding from edge
    PaddingY int      // Vertical padding from edge
}
```

### Functions

| Function | Description |
|----------|-------------|
| `DefaultOptions()` | Returns defaults (BottomRight, 0.5 opacity, 10px padding) |
| `Apply(src, wm, opts)` | Apply watermark at position |
| `ApplyFromFiles(src, wm, opts)` | Load images and apply watermark |
| `Tile(src, wm, opacity, spacing)` | Create tiled watermark pattern |
| `SaveJPEG(img, w, quality)` | Save as JPEG |
| `SavePNG(img, w)` | Save as PNG |

## Best Practices

### Watermark Image Tips

1. Use PNG format with transparency
2. Keep watermark relatively small (10-20% of image width)
3. Use white or light colors for photos
4. Add a subtle drop shadow for visibility on varied backgrounds

### Opacity Guidelines

| Opacity | Use Case |
|---------|----------|
| 0.2-0.3 | Subtle branding |
| 0.4-0.5 | Standard watermark |
| 0.6-0.8 | Prominent branding |
| 0.9-1.0 | Maximum visibility |

## Requirements

- Go 1.16 or later

## Related Packages

- [imgutils-filter](https://github.com/imgutils-org/imgutils-filter) - Image filters
- [imgutils-merge](https://github.com/imgutils-org/imgutils-merge) - Image merging
- [imgutils-sdk](https://github.com/imgutils-org/imgutils-sdk) - Unified SDK

## License

MIT License - see [LICENSE](LICENSE) for details.
