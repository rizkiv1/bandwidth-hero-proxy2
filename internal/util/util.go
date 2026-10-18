package util

import (
	"strings"
)

const (
	MinCompressLength            = 1024
	MinTransparentCompressLength = 102400
)

// ShouldCompress determines if the image should be compressed based on its type, size, and transparency.
func ShouldCompress(imageType string, size int, isTransparent bool) bool {
	if !strings.HasPrefix(imageType, "image") {
		return false
	}
	if size == 0 {
		return false
	}
	if isTransparent && size < MinCompressLength {
		return false
	}
	if !isTransparent && (strings.HasSuffix(imageType, "png") || strings.HasSuffix(imageType, "gif")) && size < MinTransparentCompressLength {
		return false
	}
	return true
}

// Pick extracts the specified keys from the given headers map.
func Pick(headers map[string][]string, properties []string) map[string]string {
	picked := make(map[string]string)

	// Create a map for fast lookup of properties (case-insensitive)
	propsMap := make(map[string]bool)
	for _, p := range properties {
		propsMap[strings.ToLower(p)] = true
	}

	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if propsMap[lowerKey] && len(values) > 0 {
			picked[lowerKey] = values[0]
		}
	}
	return picked
}
