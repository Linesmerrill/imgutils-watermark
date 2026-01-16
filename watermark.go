// Package watermark provides image watermarking utilities.
package watermark

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

// Position specifies where to place the watermark.
type Position int

const (
	// Center places the watermark in the center.
	Center Position = iota
	// TopLeft places the watermark in the top-left corner.
	TopLeft
	// TopRight places the watermark in the top-right corner.
	TopRight
	// BottomLeft places the watermark in the bottom-left corner.
	BottomLeft
	// BottomRight places the watermark in the bottom-right corner.
	BottomRight
)

// Options configures watermark placement.
type Options struct {
	Position Position
	Opacity  float64 // 0.0 to 1.0
	PaddingX int
	PaddingY int
}

// DefaultOptions returns sensible watermark defaults.
func DefaultOptions() Options {
	return Options{
		Position: BottomRight,
		Opacity:  0.5,
		PaddingX: 10,
		PaddingY: 10,
	}
}

// Apply applies a watermark image to the source image.
func Apply(src, watermark image.Image, opts Options) image.Image {
	srcBounds := src.Bounds()
	wmBounds := watermark.Bounds()

	dst := image.NewRGBA(srcBounds)
	draw.Draw(dst, srcBounds, src, srcBounds.Min, draw.Src)

	// Calculate watermark position
	var x, y int
	switch opts.Position {
	case Center:
		x = (srcBounds.Dx() - wmBounds.Dx()) / 2
		y = (srcBounds.Dy() - wmBounds.Dy()) / 2
	case TopLeft:
		x = opts.PaddingX
		y = opts.PaddingY
	case TopRight:
		x = srcBounds.Dx() - wmBounds.Dx() - opts.PaddingX
		y = opts.PaddingY
	case BottomLeft:
		x = opts.PaddingX
		y = srcBounds.Dy() - wmBounds.Dy() - opts.PaddingY
	case BottomRight:
		x = srcBounds.Dx() - wmBounds.Dx() - opts.PaddingX
		y = srcBounds.Dy() - wmBounds.Dy() - opts.PaddingY
	}

	// Apply watermark with opacity
	if opts.Opacity <= 0 {
		opts.Opacity = 0.5
	}
	if opts.Opacity > 1 {
		opts.Opacity = 1
	}

	for wy := 0; wy < wmBounds.Dy(); wy++ {
		for wx := 0; wx < wmBounds.Dx(); wx++ {
			dx := x + wx
			dy := y + wy
			if dx < 0 || dx >= srcBounds.Dx() || dy < 0 || dy >= srcBounds.Dy() {
				continue
			}

			srcColor := dst.At(dx, dy)
			wmColor := watermark.At(wmBounds.Min.X+wx, wmBounds.Min.Y+wy)

			blended := blendColors(srcColor, wmColor, opts.Opacity)
			dst.Set(dx, dy, blended)
		}
	}

	return dst
}

// blendColors blends two colors with the given opacity for the overlay.
func blendColors(base, overlay color.Color, opacity float64) color.Color {
	br, bg, bb, ba := base.RGBA()
	or, og, ob, oa := overlay.RGBA()

	// If watermark pixel is transparent, keep base
	if oa == 0 {
		return base
	}

	// Apply opacity to overlay alpha
	overlayAlpha := float64(oa) / 65535.0 * opacity

	// Blend
	r := uint8((float64(br>>8)*(1-overlayAlpha) + float64(or>>8)*overlayAlpha))
	g := uint8((float64(bg>>8)*(1-overlayAlpha) + float64(og>>8)*overlayAlpha))
	b := uint8((float64(bb>>8)*(1-overlayAlpha) + float64(ob>>8)*overlayAlpha))

	return color.RGBA{r, g, b, uint8(ba >> 8)}
}

// ApplyFromFiles loads images and applies a watermark.
func ApplyFromFiles(srcPath, watermarkPath string, opts Options) (image.Image, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	src, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, err
	}

	wmFile, err := os.Open(watermarkPath)
	if err != nil {
		return nil, err
	}
	defer wmFile.Close()

	wm, _, err := image.Decode(wmFile)
	if err != nil {
		return nil, err
	}

	return Apply(src, wm, opts), nil
}

// Tile applies a watermark in a tiled pattern across the image.
func Tile(src, watermark image.Image, opacity float64, spacing int) image.Image {
	srcBounds := src.Bounds()
	wmBounds := watermark.Bounds()

	dst := image.NewRGBA(srcBounds)
	draw.Draw(dst, srcBounds, src, srcBounds.Min, draw.Src)

	wmW := wmBounds.Dx() + spacing
	wmH := wmBounds.Dy() + spacing

	for y := 0; y < srcBounds.Dy(); y += wmH {
		for x := 0; x < srcBounds.Dx(); x += wmW {
			for wy := 0; wy < wmBounds.Dy(); wy++ {
				for wx := 0; wx < wmBounds.Dx(); wx++ {
					dx := x + wx
					dy := y + wy
					if dx >= srcBounds.Dx() || dy >= srcBounds.Dy() {
						continue
					}

					srcColor := dst.At(dx, dy)
					wmColor := watermark.At(wmBounds.Min.X+wx, wmBounds.Min.Y+wy)
					blended := blendColors(srcColor, wmColor, opacity)
					dst.Set(dx, dy, blended)
				}
			}
		}
	}

	return dst
}

// SaveJPEG saves the watermarked image as JPEG.
func SaveJPEG(img image.Image, w io.Writer, quality int) error {
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

// SavePNG saves the watermarked image as PNG.
func SavePNG(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}
