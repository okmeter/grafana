package interceptor

import (
	"net/http"
	"net/url"
)

// WithRequestHeader sets additional (or overrides existing) header value to request
func WithRequestHeader(key, value string) Interceptor {
	return func(next Doer) Doer {
		return func(req *http.Request) (*http.Response, error) {
			req.Header.Set(key, value)
			return next(req)
		}
	}
}

// WithRequestCookie adds cookie parameter to request
func WithRequestCookie(cookie *http.Cookie) Interceptor {
	return func(next Doer) Doer {
		return func(req *http.Request) (*http.Response, error) {
			req.AddCookie(cookie)
			return next(req)
		}
	}
}

// WithRequestQueryParams sets query parameters to request
func WithRequestQueryParams(params url.Values) Interceptor {
	return func(next Doer) Doer {
		return func(req *http.Request) (*http.Response, error) {
			req.URL.RawQuery = params.Encode()
			return next(req)
		}
	}
}
