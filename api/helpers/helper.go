package helpers

import (
	"os"
	"strings"
)

// Makes Sure that the URL has a prefix of http or https
func EnsuredPrefixHTTP(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

func IsDifferentDomain(url string) bool {
	domain := os.Getenv("DOMAIN")

	if url == domain {
		return false
	}

	cleanURL := strings.TrimPrefix(url, "http://")
	cleanURL = strings.TrimPrefix(cleanURL, "https://")
	cleanURL = strings.TrimPrefix(cleanURL,"www.")
	cleanURL = strings.Split(cleanURL, "/")[0]

	// Compare the cleaned URL with the domain
	return cleanURL != domain
}

// http://example.com			example.com
// https://example.com			example.com
// http://www.example.com		example.com
// https://www.example.com		example.com
// www.example.com				example.com
// example.com/about			example.com
// http://example.com/about		example.com
// http://example.com:8080		example.com:8080