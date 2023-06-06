// Package uadeco user-agent decorator package
// should be put on internal/ of the project's directory
// serviceName should be changed base on the project name/directory
package uadeco

import (
	"io"
	"net/http"
	"runtime/debug"
	"strings"
)

// nolint: gochecknoglobals
var serviceName = `uadeco`

// nolint: gochecknoglobals
var userAgent = `uadeco/0.0`

// nolint: gochecknoglobals
var defaultHeaders = map[string]string{
	`User-Agent`: userAgent,
}

// OriginalHttpTransport is backup original http transport
// nolint: gochecknoglobals
var OriginalHttpTransport = http.DefaultTransport

// OriginalHttpClient is backup original http client
// nolint: gochecknoglobals
var OriginalHttpClient = http.DefaultClient

// DecoratedTransport transport with custom user-agent for http.DefaultTransport
// nolint: gochecknoglobals
var DecoratedTransport = &Transport{
	headers: defaultHeaders,
}

// DecoratedClient client with custom user-agent for http.DefaultClient
// nolint: gochecknoglobals
var DecoratedClient = &http.Client{
	Transport: DecoratedTransport,
}

// InitServiceName must be called on init of project's main or TestMain
// for example:
//
//	uadec.InitServiceName("service1")
func InitServiceName(name string) {
	defer ReplaceDefaultTransport()
	serviceName = name
	userAgent = serviceName + `/`
	if info, ok := debug.ReadBuildInfo(); ok {
		var rev, time string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision": // eg. vcs.revision=258a9cbd34e46ea419ffa6cc4d3f37f3e2ca4277
				rev = setting.Value
			case "vcs.time": // eg. vcs.time=2022-10-21T10:18:32Z
				time = setting.Value
			}
		}
		{ // build suffix
			ver := ``
			addSuffix := func(s string) {
				if s == `` {
					return
				}
				if ver != `` {
					ver += `-`
				}
				ver += s
			}
			time = strings.ReplaceAll(time, `:`, ``)
			time = strings.ReplaceAll(time, `-`, ``)
			addSuffix(time)
			addSuffix(rev)
			addSuffix(info.GoVersion)
			userAgent += ver
			return
		}
	}
	userAgent += `0.0-dev`
}

// SetUserAgent used if need to set user agent manually
func SetUserAgent(r *http.Request) {
	r.Header.Set("User-Agent", userAgent)
}

// NewHttpRequest a replacement for uadec.NewHttpRequest with decorated user agent
// nolint: wrapcheck
func NewHttpRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if req != nil {
		SetUserAgent(req)
	}
	return req, err
}

// HttpGet a replacement for http.Get with decorated user agent
// use this when http.DefaultTransport being replaced by other package
// nolint: wrapcheck
func HttpGet(url string) (*http.Response, error) {
	req, err := NewHttpRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// Transport class that stores haeder and original roundtripper
type Transport struct {
	headers map[string]string
}

// RoundTrip add headers and use base roundtripper
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return OriginalHttpTransport.RoundTrip(req) // nolint: wrapcheck
}

// ReplaceDefaultTransport replace http.DefaultTransport with DecoratedTransport
//
//	also replaces http.DefaultClient with Decorated
func ReplaceDefaultTransport() {
	defaultHeaders[`User-Agent`] = userAgent
	http.DefaultTransport = DecoratedTransport
	http.DefaultClient = DecoratedClient
}
