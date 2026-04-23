package logic

import (
	"net/url"
	"strings"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
)

// extractDomain parses rawURL and returns the hostname with www. stripped.
// Returns a ValidationError if the URL is malformed or missing scheme/host.
func extractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", commonerrors.New("url must be a valid absolute URL (e.g. https://example.com/path)")
	}
	host := strings.TrimPrefix(u.Hostname(), "www.")
	return host, nil
}

// computeReadingTime returns ceil(wordCount / 200), minimum 1.
func computeReadingTime(body string) int {
	words := len(strings.Fields(body))
	if words == 0 {
		return 1
	}
	return (words + 199) / 200
}
