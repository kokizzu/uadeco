package uadeco

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

func TestUserAgentIncludesServiceAndGoVersion(t *testing.T) {
	ua := configureForTest(t, "uadeco")
	if !strings.HasPrefix(ua, "uadeco/") {
		t.Fatalf("user agent should include service prefix: %q", ua)
	}
	if !strings.Contains(ua, runtime.Version()) {
		t.Fatalf("user agent should include current Go version %q: %q", runtime.Version(), ua)
	}
}

func TestNewHTTPRequestSetsUserAgent(t *testing.T) {
	ua := configureForTest(t, "uadeco")
	req, err := NewHttpRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewHttpRequest error: %v", err)
	}
	if got := req.Header.Get("User-Agent"); got != ua {
		t.Fatalf("User-Agent want %q got %q", ua, got)
	}
}

func TestDecoratedClientTransport(t *testing.T) {
	ua := configureForTest(t, "uadeco")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, r.Header.Get("User-Agent"))
	}))
	t.Cleanup(server.Close)

	tests := []struct {
		name string
		get  func() (*http.Response, error)
	}{
		{
			name: "HttpGet",
			get: func() (*http.Response, error) {
				return HttpGet(server.URL)
			},
		},
		{
			name: "decorated-client",
			get: func() (*http.Response, error) {
				return DecoratedClient.Get(server.URL)
			},
		},
		{
			name: "default-client",
			get: func() (*http.Response, error) {
				return http.Get(server.URL)
			},
		},
		{
			name: "custom-client",
			get: func() (*http.Response, error) {
				client := &http.Client{}
				return client.Get(server.URL)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := tc.get()
			if err != nil {
				t.Fatalf("request error: %v", err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Fatalf("closing response body: %v", err)
				}
			}()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("reading response body: %v", err)
			}
			if got := string(body); got != ua {
				t.Fatalf("User-Agent want %q got %q", ua, got)
			}
		})
	}
}

func configureForTest(t *testing.T, name string) string {
	t.Helper()

	originalServiceName := serviceName
	originalUserAgent := userAgent
	originalHeader := defaultHeaders["User-Agent"]
	originalTransport := http.DefaultTransport
	originalClient := http.DefaultClient

	t.Cleanup(func() {
		serviceName = originalServiceName
		userAgent = originalUserAgent
		defaultHeaders["User-Agent"] = originalHeader
		http.DefaultTransport = originalTransport
		http.DefaultClient = originalClient
	})

	InitServiceName(name)
	return userAgent
}
