package image

import (
	"bytes"
	"image"
	"io"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// Processor handles image processing operations
type Processor struct {
	quality int // WebP quality (0-100)
}

// NewProcessor creates a new image processor with the given quality setting
func NewProcessor(quality int) *Processor {
	if quality < 0 {
		quality = 0
	}
	if quality > 100 {
		quality = 100
	}
	return &Processor{quality: quality}
}

// ProcessResult contains the processed image data and metadata
type ProcessResult struct {
	Data   []byte
	Width  int
	Height int
}

// DecodeImage decodes an image from a reader
func (p *Processor) DecodeImage(r io.Reader) (image.Image, error) {
	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// GetDimensions returns the width and height of an image
func (p *Processor) GetDimensions(img image.Image) (width, height int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// ResizeToWidth resizes an image to the specified max width while maintaining aspect ratio
// If the image is smaller than maxWidth, it returns the original image unchanged
func (p *Processor) ResizeToWidth(img image.Image, maxWidth int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()

	if width <= maxWidth {
		return img
	}

	return imaging.Resize(img, maxWidth, 0, imaging.Lanczos)
}

// EncodeToWebP encodes an image to WebP format with the processor's quality setting
func (p *Processor) EncodeToWebP(img image.Image) ([]byte, error) {
	var buf bytes.Buffer

	options := &webp.Options{
		Lossless: false,
		Quality:  float32(p.quality),
	}

	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ProcessImage processes an image: decodes, optionally resizes, and encodes to WebP
// If maxWidth is 0, no resizing is performed
func (p *Processor) ProcessImage(r io.Reader, maxWidth int) (*ProcessResult, error) {
	img, err := p.DecodeImage(r)
	if err != nil {
		return nil, err
	}

	if maxWidth > 0 {
		img = p.ResizeToWidth(img, maxWidth)
	}

	width, height := p.GetDimensions(img)

	data, err := p.EncodeToWebP(img)
	if err != nil {
		return nil, err
	}

	return &ProcessResult{
		Data:   data,
		Width:  width,
		Height: height,
	}, nil
}

// GenerateThumbnails generates thumbnails at specified widths
// Returns a map of suffix -> ProcessResult
func (p *Processor) GenerateThumbnails(img image.Image, sizes map[string]int) (map[string]*ProcessResult, error) {
	results := make(map[string]*ProcessResult)

	for suffix, maxWidth := range sizes {
		resized := p.ResizeToWidth(img, maxWidth)
		width, height := p.GetDimensions(resized)

		data, err := p.EncodeToWebP(resized)
		if err != nil {
			return nil, err
		}

		results[suffix] = &ProcessResult{
			Data:   data,
			Width:  width,
			Height: height,
		}
	}

	return results, nil
}
