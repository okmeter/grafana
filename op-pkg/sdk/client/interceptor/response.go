package interceptor

import "net/http"

// WithResponseCodeCustomError sets custom error on given response code
func WithResponseCodeCustomError(responseCode int, customErr error) Interceptor {
	return func(next Doer) Doer {
		return func(req *http.Request) (*http.Response, error) {
			res, err := next(req)
			if err != nil {
				return nil, err
			}
			if res.StatusCode == responseCode {
				return nil, customErr
			}
			return res, nil
		}
	}
}
