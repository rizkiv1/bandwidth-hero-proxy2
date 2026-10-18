package proxy

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bandwidth-hero-proxy/internal/compress"
	"bandwidth-hero-proxy/internal/util"
)

const DefaultQuality = 40

var forwardedHeaders = []string{
	"cookie", "dnt", "referer", "user-agent", "accept", "accept-language", "accept-encoding",
}

var bmiRegex = regexp.MustCompile(`(?i)http://1\.1\.\d\.\d/bmi/(https?://)?`)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

// Handler serves as the main HTTP handler for the compression proxy
func Handler(w http.ResponseWriter, r *http.Request) {
	// Root endpoint logic from Netlify config (index handler)
	if r.URL.Path != "/api/index" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("bandwidth-hero-proxy"))
		return
	}

	// Sometimes the URL comes encoded as JSON string?
	var parsedUrl string
	if err := json.Unmarshal([]byte(urlParam), &parsedUrl); err == nil {
		urlParam = parsedUrl
	}

	// If it's an array of URLs joined (sometimes happens with bad query building)
	// We just take the first one or fix it. The JS implementation has:
	// Array.isArray(r) && (r = r.join("&url="))
	// In Go, if multiple `url` params are passed, `r.URL.Query().Get("url")` returns the first one.
	// But let's check if the raw query has multiple `url=`
	if len(r.URL.Query()["url"]) > 1 {
		urlParam = strings.Join(r.URL.Query()["url"], "&url=")
	}

	urlParam = bmiRegex.ReplaceAllString(urlParam, "http://")

	useWebp := r.URL.Query().Get("jpeg") == "0"
	if r.URL.Query().Get("jpeg") == "" {
		useWebp = true // Default is true, !s in js where s is jpeg
	}

	grayscale := r.URL.Query().Get("bw") != "0" && r.URL.Query().Get("bw") != ""

	qualityStr := r.URL.Query().Get("l")
	quality, err := strconv.Atoi(qualityStr)
	if err != nil || quality <= 0 {
		quality = DefaultQuality
	}

	// Fetch image
	req, err := http.NewRequestWithContext(r.Context(), "GET", urlParam, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy allowed headers
	pickedHeaders := util.Pick(r.Header, forwardedHeaders)
	for k, v := range pickedHeaders {
		req.Header.Set(k, v)
	}

	// X-Forwarded-For
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		// Use RemoteAddr
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			xff = host
		} else {
			xff = r.RemoteAddr
		}
	}
	req.Header.Set("X-Forwarded-For", xff)

	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		w.WriteHeader(resp.StatusCode)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	originalSize := len(bodyBytes)

	// Check if should compress
	if !util.ShouldCompress(contentType, originalSize, useWebp) {
		log.Printf("Bypassing... Size: %d\n", originalSize)
		// Send original
		w.Header().Set("Content-Encoding", "identity")
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.Header().Set("Content-Length", strconv.Itoa(originalSize))
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
		return
	}

	// Compress
	result, err := compress.Compress(bodyBytes, useWebp, grayscale, quality, originalSize)
	if err != nil {
		log.Printf("Conversion failed: %s, Error: %v\n", urlParam, err)
		// Fallback to original image if compression fails
		w.Header().Set("Content-Encoding", "identity")
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.Header().Set("Content-Length", strconv.Itoa(originalSize))
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
		return
	}

	log.Printf("From %d, Saved: %.2f%%\n", originalSize, float64(originalSize-len(result.Output))/float64(originalSize)*100)

	// Send compressed
	w.Header().Set("Content-Encoding", "identity")
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	for k, v := range result.Headers {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(result.Output)))
	w.WriteHeader(http.StatusOK)
	w.Write(result.Output)
}
