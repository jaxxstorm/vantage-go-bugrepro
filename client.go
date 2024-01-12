package main

import (
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	vantagev2 "github.com/vantage-sh/vantage-go/vantagev2/vantage"
)

const userAgent = "tf-provider-vantage"

type Client struct {
	V2   *vantagev2.Vantage
	Auth runtime.ClientAuthInfoWriter
}

func NewClient(host, token string) (*Client, error) {
	parsedURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	v2Cfg := vantagev2.DefaultTransportConfig()
	v2Cfg.WithHost(parsedURL.Host)
	v2Cfg.WithSchemes([]string{parsedURL.Scheme})
	transportv2 := httptransport.New(v2Cfg.Host, v2Cfg.BasePath, v2Cfg.Schemes)
	transportv2.Transport = userAgentTripper(transportv2.Transport, userAgent)
	v2 := vantagev2.New(transportv2, strfmt.Default)

	bearerTokenAuth := httptransport.BearerToken(token)
	return &Client{
		V2:   v2,
		Auth: bearerTokenAuth,
	}, nil
}

func userAgentTripper(inner http.RoundTripper, userAgent string) http.RoundTripper {
	version := "unknown"
	modified := false
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				version = s.Value[:7]
			case "vcs.modified":
				modified = s.Value == "true"
			}
		}
	}
	agent := userAgent + "/" + version
	if modified {
		agent = agent + "+"
	}
	return &roundtripper{
		inner: inner,
		agent: agent,
	}
}

type roundtripper struct {
	inner http.RoundTripper
	agent string
}

func (ug *roundtripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", ug.agent)
	return ug.inner.RoundTrip(r)
}
