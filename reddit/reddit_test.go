package reddit

import (
	"net/http"
	"net/url"
	"testing"
)

func TestUserAgent(t *testing.T) {
	s1 := UserAgent("linux", "tl.foo.bar", "v0.0.1", "/u/anon")
	s2 := "linux:tl.foo.bar:v0.0.1 (by /u/anon)"
	if s1 != s2 {
		t.Errorf("UserAgent was %s instead of %s", s1, s2)
	}
}

func TestNewClient(t *testing.T) {
	c := &http.Client{}
	client := NewClient(c)
	if client.client != c {
		t.Errorf("http.Client was not set")
	}
}

func TestNewClientDefaults(t *testing.T) {
	client := NewClient(nil)
	if client.client != http.DefaultClient {
		t.Errorf("Defautl http.Client was not set")
	}
	if b := client.BaseURL.String(); b != defaultBaseURL {
		t.Errorf("defaultBaseURL was %s instead of defaultBaseURL(%s)",
			b, defaultBaseURL)
	}
	if client.UserAgent != defaultUserAgent {
		t.Errorf("UserAgent was %s instead of defaultUserAgent(%s)",
			client.UserAgent, defaultUserAgent)
	}
}

func TestClientNewRequest(t *testing.T) {
	client := NewClient(nil)
	bu, _ := url.Parse("http://whatever.tld:8080")
	client.BaseURL = bu
	ua := "foobar"
	client.UserAgent = ua
	req, err := client.NewRequest("GET", "/some/path.json", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "GET" {
		t.Errorf("Method was '%s' instead of GET", req.Method)
	}
	if hua := req.Header.Get("User-Agent"); hua != ua {
		t.Errorf("UserAgent Header was '%s' instead of '%s'", hua, ua)
	}
	fullURL := "http://whatever.tld:8080/some/path.json?raw_json=1"
	if s := req.URL.String(); s != fullURL {
		t.Errorf("Response URL was '%s' instead of '%s'", s, fullURL)
	}
}

func TestWrapHTTPTransport(t *testing.T) {
	transport := &http.Transport{}
	rt := WrapHTTPTransport("asdf", transport)
	if rt.RoundTripper != transport {
		t.Errorf("WrapHTTPTransport RoundTripper is %s instead of %s",
			rt.RoundTripper, transport)
	}
}

func TestWrapHTTPTransportRTFallback(t *testing.T) {
	rt := WrapHTTPTransport("asdf", nil)
	if rt.RoundTripper != http.DefaultTransport {
		t.Errorf("WrapHTTPTransport has not set the http.DefaultTransport as fallback")
	}
}

func TestNewRequestJSONHeader(t *testing.T) {
	client := NewClient(nil)
	req, err := client.NewRequest("GET", "/test?foo=bar", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	if s := q.Get("foo"); s != "bar" {
		t.Errorf("URL query foo was '%s' instead of 'bar'", s)
	}
	if q.Get("raw_json") != "1" {
		t.Error("Client requests don't set the raw_json header")
	}
}
