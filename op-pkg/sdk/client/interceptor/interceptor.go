package interceptor

import "net/http"

// Doer does an HTTP request
type Doer func(*http.Request) (*http.Response, error)

// Interceptor is a Doer decorator useful modify request/response or to extract some data
type Interceptor func(next Doer) Doer
