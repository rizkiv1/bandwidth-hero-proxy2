package compress

import (
	"fmt"
	"github.com/h2non/bimg"
)

type Result struct {
	Headers map[string]string
	Output  []byte
}

// Compress scales and converts the image using bimg (libvips)
func Compress(imageBuffer []byte, useWebp bool, grayscale bool, quality int, originalSize int) (*Result, error) {
	img := bimg.NewImage(imageBuffer)

	options := bimg.Options{
		Quality: quality,
	}

	if useWebp {
		options.Type = bimg.WEBP
	} else {
		options.Type = bimg.JPEG
		options.Interlace = true // progressive: true
	}

	if grayscale {
		// BIMG doesn't have a direct "grayscale" boolean in options, but it has interpretations
		options.Interpretation = bimg.InterpretationBW
	}

	newImage, err := img.Process(options)
	if err != nil {
		return nil, fmt.Errorf("compression failed: %w", err)
	}

	newSize := len(newImage)

	formatStr := "jpeg"
	if useWebp {
		formatStr = "webp"
	}

	headers := map[string]string{
		"content-type":    fmt.Sprintf("image/%s", formatStr),
		"content-length":  fmt.Sprintf("%d", newSize),
		"x-original-size": fmt.Sprintf("%d", originalSize),
		"x-bytes-saved":   fmt.Sprintf("%d", originalSize-newSize),
	}

	return &Result{
		Headers: headers,
		Output:  newImage,
	}, nil
}
