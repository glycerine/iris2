package httptest

import (
	"bufio"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/go-iris2/iris2"
)

type HTTPTest struct {
	t   *testing.T
	srv *iris2.Server
	l   *fasthttputil.InmemoryListener
}

// New Prepares and returns a new test framework based on the app
func New(app *iris2.Server, t *testing.T) *HTTPTest {

	ht := &HTTPTest{
		t:   t,
		srv: app,
		l:   fasthttputil.NewInmemoryListener(),
	}
	go func() {
		err := ht.srv.Serve(ht.l)
		if err != nil {
			panic(err)
		}
	}()
	return ht
}

func (h *HTTPTest) Get(url string) *Expecter {
	return h.Request("GET", url)
}

func (h *HTTPTest) Post(url string) *Expecter {
	return h.Request("GET", url)
}

func (h *HTTPTest) Request(method, url string) *Expecter {
	if !strings.Contains(url, "http://") {
		url = "http://localhost" + url
	}
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("User-Agent", "HTTPTest/v1.0")

	dump, err := httputil.DumpRequestOut(req, true)
	require.Nil(h.t, err)

	c, err := h.l.Dial()
	require.Nil(h.t, err)

	_, err = c.Write(dump)
	require.Nil(h.t, err)

	br := bufio.NewReader(c)
	var resp fasthttp.Response
	err = resp.Read(br)
	require.Nil(h.t, err)

	var body []byte
	if string(resp.Header.Peek("Content-Encoding")) == "gzip" {
		body, err = resp.BodyGunzip()
		require.Nil(h.t, err)
	} else {
		body = resp.Body()
	}

	return &Expecter{
		body: string(body),
		t:    h.t,
		resp: &resp,
	}
}
