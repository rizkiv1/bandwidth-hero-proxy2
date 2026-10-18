package util

import (
	"testing"
)

func TestShouldCompress(t *testing.T) {
	tests := []struct {
		name          string
		imageType     string
		size          int
		isTransparent bool
		expected      bool
	}{
		{"Not an image", "text/html", 2000, false, false},
		{"Zero size", "image/jpeg", 0, false, false},
		{"Transparent too small", "image/png", 500, true, false},
		{"Transparent big enough", "image/png", 2000, true, true},
		{"Non-transparent PNG too small", "image/png", 50000, false, false},
		{"Non-transparent PNG big enough", "image/png", 150000, false, true},
		{"JPEG always compress if > 0", "image/jpeg", 500, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldCompress(tt.imageType, tt.size, tt.isTransparent); got != tt.expected {
				t.Errorf("ShouldCompress() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPick(t *testing.T) {
	headers := map[string][]string{
		"Cookie":       {"test=1"},
		"Dnt":          {"1"},
		"RandomHeader": {"ignore"},
	}

	properties := []string{"cookie", "dnt", "missing"}
	picked := Pick(headers, properties)

	if len(picked) != 2 {
		t.Errorf("Expected 2 picked headers, got %d", len(picked))
	}
	if picked["cookie"] != "test=1" {
		t.Errorf("Expected test=1 for cookie")
	}
	if picked["dnt"] != "1" {
		t.Errorf("Expected 1 for dnt")
	}
}
