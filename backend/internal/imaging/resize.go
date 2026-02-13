package imaging

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

// ResizeImage resizes an image to fit within maxWidth x maxHeight while maintaining aspect ratio.
// Returns the resized image as JPEG bytes. Returns original bytes if not a supported image format.
func ResizeImage(data []byte, maxWidth, maxHeight int) ([]byte, string, error) {
	src, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		// Not a decodable image, return original
		return data, "", err
	}

	bounds := src.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	// No resize needed if already within bounds
	if origW <= maxWidth && origH <= maxHeight {
		return data, format, nil
	}

	// Calculate new dimensions maintaining aspect ratio
	newW, newH := fitDimensions(origW, origH, maxWidth, maxHeight)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	var buf bytes.Buffer
	outFormat := format
	switch format {
	case "png":
		err = png.Encode(&buf, dst)
	default:
		outFormat = "jpeg"
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85})
	}
	if err != nil {
		return data, format, err
	}

	return buf.Bytes(), outFormat, nil
}

// GenerateThumbnail creates a square center-cropped thumbnail of the given size.
func GenerateThumbnail(data []byte, size int) ([]byte, string, error) {
	src, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	bounds := src.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	// Center crop to square
	cropSize := origW
	if origH < cropSize {
		cropSize = origH
	}
	x0 := (origW - cropSize) / 2
	y0 := (origH - cropSize) / 2

	cropped := image.NewRGBA(image.Rect(0, 0, cropSize, cropSize))
	draw.Copy(cropped, image.Point{}, src, image.Rect(x0, y0, x0+cropSize, y0+cropSize), draw.Over, nil)

	// Resize to target size
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), cropped, cropped.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	outFormat := format
	switch format {
	case "png":
		err = png.Encode(&buf, dst)
	default:
		outFormat = "jpeg"
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80})
	}
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), outFormat, nil
}

// ResizeFromReader is a convenience wrapper that reads from an io.Reader.
func ResizeFromReader(r io.Reader, maxWidth, maxHeight int) ([]byte, string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, "", err
	}
	return ResizeImage(data, maxWidth, maxHeight)
}

func fitDimensions(origW, origH, maxW, maxH int) (int, int) {
	ratioW := float64(maxW) / float64(origW)
	ratioH := float64(maxH) / float64(origH)
	ratio := ratioW
	if ratioH < ratio {
		ratio = ratioH
	}
	return int(float64(origW) * ratio), int(float64(origH) * ratio)
}
