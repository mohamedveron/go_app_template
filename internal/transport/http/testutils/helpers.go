package testutils

import (
	"net/url"
	"strings"
)

// EncodePathForURL encodes a file path for use in URLs
func EncodePathForURL(filePath string) string {
	// Remove leading slash and encode each component
	path := strings.TrimPrefix(filePath, "/")
	parts := strings.Split(path, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

// TestAPIPath returns the v1 API path for a given endpoint
func TestAPIPath(endpoint string) string {
	// Ensure endpoint starts with /
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	return "/api/v1" + endpoint
}
