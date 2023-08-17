package client

import (
	"net/http"
)

// Option defines an option for a Client
type Option func(*Client)

// OptionHTTPClient sets another http.Client
func OptionHTTPClient(httpClient *http.Client) func(*Client) {
	return func(c *Client) { c.httpClient = httpClient }
}

// OptionUserAgent rewrites "User-Agent" header value
func OptionUserAgent(userAgent string) func(*Client) {
	return OptionHeader("User-Agent", userAgent)
}

// OptionHeader adds or sets existing header value
func OptionHeader(key, value string) func(*Client) {
	return func(c *Client) { c.headers[key] = value }
}

// OptionTransport sets another http.Client transport
func OptionTransport(transport http.RoundTripper) func(*Client) {
	return func(c *Client) { c.httpClient.Transport = transport }
}
