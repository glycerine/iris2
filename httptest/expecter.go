package httptest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

type Expecter struct {
	body string
	t    *testing.T
	resp *fasthttp.Response
}

func (e *Expecter) Status(code int) *Expecter {
	require.Equal(e.t, code, e.resp.Header.StatusCode(), "unexpected status code")
	return e
}

func (e *Expecter) Body(body string) *Expecter {
	require.Equal(e.t, body, e.body, "unexpected body")
	return e
}
