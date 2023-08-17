package roundtripper

import (
	"bytes"
	"io"
	"net/http"

	"github.com/grafana/grafana/op-pkg/sdk/middleware"

	"github.com/grafana/grafana/pkg/infra/log"
)

func NewLoggingRoundTripper(transport http.RoundTripper, component string) http.RoundTripper {
	return roundTripper{
		internal: transport,
		logger:   log.New(component),
	}
}

type roundTripper struct {
	internal http.RoundTripper
	logger   log.Logger
}

func (rt roundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	var (
		querier = middleware.GetQuerier(req.Context())
	)
	switch req.Method {
	case http.MethodPost, http.MethodPut:
		r1, r2, _ := drainBody(req.Body)
		payload, _ := readBody(r1)
		req.Body = r2
		rt.logger.Debug("roundtrip", "querier", querier, "method", req.Method, "url", req.URL, "payload", payload)
	default:
		rt.logger.Debug("roundtrip", "querier", querier, "method", req.Method, "url", req.URL)
	}
	return rt.internal.RoundTrip(req)
}

func readBody(b io.ReadCloser) (s string, err error) {
	var buff bytes.Buffer
	_, err = buff.ReadFrom(b)
	if err != nil {
		return
	}
	s = buff.String()
	return
}

// drainBody copied from https://go.dev/src/net/http/httputil/dump.go
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
