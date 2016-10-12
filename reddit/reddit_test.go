package reddit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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

func TestNewRequestWithBody(t *testing.T) {
	body := map[string]string{"foo": "bar"}
	req, err := NewClient(nil).NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	var (
		is     = string(reqBody)
		should = "{\"foo\":\"bar\"}\n"
	)
	if is != should {
		t.Errorf("Body was '%#v' instead of '%#v'", is, should)
	}
}

func TestNewRequestWithBrokenBody(t *testing.T) {
	var body = jsonEncodingBreaker{}
	req, err := NewClient(nil).NewRequest("POST", "/", body)
	if err == nil {
		r, _ := ioutil.ReadAll(req.Body)
		t.Errorf("No error for broken body: %s\n", r)
	}
}

type jsonEncodingBreaker struct{}

func (jsonEncodingBreaker) MarshalJSON() ([]byte, error) {
	return nil, errors.New("ERROR")
}

func TestClientDoClientErr(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, errIs := NewClient(&http.Client{
		Transport: brokenRoundTripper{Error: errors.New("BREAK")},
	}).Do(req, nil)
	if errIs == nil {
		t.Error("Returned no error")
	}
}

func TestClientDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `Everthing is fine`)
	}))
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = NewClient(nil).Do(req, nil)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClientDoParseJSONResp(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"foo":"bar"}`)
	}))
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	var body map[string]string
	_, err = NewClient(nil).Do(req, &body)
	if err != nil {
		t.Error(err)
		return
	}
	if body["foo"] != "bar" {
		t.Errorf(`Parsed '%#v' instead of '{"foo":"bar"}'`, body)
	}

}

type brokenRoundTripper struct{ Error error }

func (t brokenRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, t.Error
}

func TestClientDoRateLimitErr(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ratelimit-Used", "10")
		w.Header().Set("X-Ratelimit-Remaining", "20")
		w.Header().Set("X-Ratelimit-Reset", "30")
		w.WriteHeader(429)
	}))
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = NewClient(nil).Do(req, nil)
	rateErr, ok := err.(*RateLimitError)
	if !ok {
		t.Error("No RateLimitError returned:", err.Error())
		return
	}
	e := &RateLimitError{
		Used:      10,
		Remaining: 20,
		Reset:     30,
	}
	if !reflect.DeepEqual(rateErr, e) {
		t.Error("RateLimitError was '%#v' instead of '%#v'", rateErr, e)
	}
}

func TestClientDoRateLimitErrUnparsable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ratelimit-Used", "10")
		w.Header().Set("X-Ratelimit-Remaining", "20")
		w.Header().Set("X-Ratelimit-Reset", "30")
		w.WriteHeader(429)
	}))
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = NewClient(nil).Do(req, nil)
	rateErr, ok := err.(*RateLimitError)
	if !ok {
		t.Error("No RateLimitError returned:", err.Error())
		return
	}
	e := &RateLimitError{
		Used:      10,
		Remaining: 20,
		Reset:     30,
	}
	if !reflect.DeepEqual(rateErr, e) {
		t.Error("RateLimitError was '%#v' instead of '%#v'", rateErr, e)
	}
}

func TestClientCheckResponseOK(t *testing.T) {
	err := CheckResponse(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`OK`)),
	})
	if err != nil {
		t.Error("Returned error for valid response: ", err)
	}
}

func TestClientCheckResponseJSONError(t *testing.T) {
	err := CheckResponse(&http.Response{
		StatusCode: http.StatusNotFound,
		Body: ioutil.NopCloser(strings.NewReader(`{
			"message": "Not found",
			"error": 404
		}`)),
	})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Error("Returned error wrong error: ", err)
	}
	shouldErr := APIError{
		Message:   "Not found",
		ErrorCode: http.StatusNotFound,
	}
	if !reflect.DeepEqual(*apiErr, shouldErr) {
		t.Errorf("Returned '%#v' instead of '%#v'", *apiErr, shouldErr)
	}

}

func TestClientCheckResponseBrokenJSONError(t *testing.T) {
	const s = `{
		broken json
	}`
	err := CheckResponse(&http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(strings.NewReader(s)),
	})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Error("Returned error wrong error: ", err)
	}
	shouldErr := APIError{
		Message:   s,
		ErrorCode: http.StatusNotFound,
	}
	if !reflect.DeepEqual(*apiErr, shouldErr) {
		t.Errorf("Returned '%#v' instead of '%#v'", *apiErr, shouldErr)
	}

}
