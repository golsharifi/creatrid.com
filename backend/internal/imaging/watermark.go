package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

// ApplyWatermark composites the watermark image onto the bottom-right corner
// of the base image at the given opacity (0..1). The mark is scaled to ~22% of
// the base width. Returns the watermarked image in the base image's format
// (PNG stays PNG, everything else re-encodes as JPEG).
func ApplyWatermark(base, mark []byte, opacity float64) ([]byte, string, error) {
	src, format, err := image.Decode(bytes.NewReader(base))
	if err != nil {
		return nil, "", err
	}
	wm, _, err := image.Decode(bytes.NewReader(mark))
	if err != nil {
		return nil, "", err
	}

	bounds := src.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, src, bounds.Min, draw.Src)

	// Scale the mark to ~22% of base width, preserving its aspect ratio.
	markW := bounds.Dx() * 22 / 100
	if markW < 24 {
		markW = 24
	}
	wmBounds := wm.Bounds()
	markH := markW * wmBounds.Dy() / wmBounds.Dx()
	scaled := image.NewRGBA(image.Rect(0, 0, markW, markH))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), wm, wmBounds, draw.Over, nil)

	// Bottom-right with a margin of 3% of the base width.
	margin := bounds.Dx() * 3 / 100
	pos := image.Rect(
		bounds.Max.X-markW-margin,
		bounds.Max.Y-markH-margin,
		bounds.Max.X-margin,
		bounds.Max.Y-margin,
	)

	if opacity <= 0 || opacity > 1 {
		opacity = 0.45
	}
	mask := image.NewUniform(color.Alpha{A: uint8(opacity * 255)})
	draw.DrawMask(canvas, pos, scaled, image.Point{}, mask, image.Point{}, draw.Over)

	var buf bytes.Buffer
	outFormat := "jpeg"
	if format == "png" {
		outFormat = "png"
		err = png.Encode(&buf, canvas)
	} else {
		err = jpeg.Encode(&buf, canvas, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), outFormat, nil
}
