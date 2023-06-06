package uadeco

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommitHash(t *testing.T) {
	InitServiceName(`uadeco`)
	if !strings.Contains(userAgent, `go1.19`) &&
		!strings.Contains(userAgent, `go1.20`) {
		t.Errorf(`Test run using outdated go version %v`, userAgent)
	}
	fmt.Println(userAgent)
}

type FakeClient struct {
	RoundTripFunc func(r *http.Request) (*http.Response, error)
}

func (s *FakeClient) RoundTrip(r *http.Request) (*http.Response, error) {
	return s.RoundTripFunc(r)
}

func TestUADeco(t *testing.T) {
	var rt http.RoundTripper = &FakeClient{
		RoundTripFunc: func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, userAgent, r.Header.Get("User-Agent"))
			return &http.Response{}, nil
		},
	}
	req, err := NewHttpRequest("GET", "/", nil)
	assert.Nil(t, err)
	client := http.Client{
		Transport: rt,
	}
	_, _ = client.Do(req)
}

func TestDecoratedClientTransport(t *testing.T) {
	const testHost = `localhost:65533`

	go func() {
		http.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, r.Header.Get(`user-agent`))
		})
		_ = http.ListenAndServe(testHost, nil)
	}()

	compareUA := func(t *testing.T, resp *http.Response, err error) {
		if assert.Nil(t, err) {
			ua, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, userAgent, string(ua))
		}
	}

	t.Run(`HttpGet`, func(t *testing.T) {
		for z := 0; z < 10; z++ {
			resp, err := HttpGet(`http://` + testHost)
			if err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			compareUA(t, resp, err)
		}
	})

	t.Run(`decorated-client`, func(t *testing.T) {
		client := DecoratedClient
		resp, err := client.Get(`http://` + testHost)
		compareUA(t, resp, err)
	})

	t.Run(`default-client`, func(t *testing.T) {
		resp, err := http.Get(`http://` + testHost)
		compareUA(t, resp, err)
	})

	t.Run(`custom-client`, func(t *testing.T) {
		client := &http.Client{}
		resp, err := client.Get(`http://` + testHost)
		compareUA(t, resp, err)
	})
}
