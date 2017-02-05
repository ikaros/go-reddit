package reddit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"
)

// Client handles communication with the reddit API.
type Client struct {
	client    *http.Client
	UserAgent string
	BaseURL   *url.URL

	common service

	Account         *AccountService
	Captcha         *CaptchaService
	Flair           *FlairService
	Gold            *GoldService
	LinksComments   *LinksCommentsService
	Listings        *ListingsService
	LiveThreads     *LiveThreadsService
	PrivateMessages *PrivateMessagesService
	Misc            *MiscService
	Moderation      *ModerationService
	Multis          *MultisService
	Search          *SearchService
	Subreddits      *SubredditsService
	Users           *UsersService
	Wiki            *WikiService

	rateLimitMu sync.Mutex
	rateLimit   *rateLimit
}

// Semantic Version
var Version = "v0.0.1"

var (
	defaultUserAgent = UserAgent(runtime.GOOS, "com.github.reddit-go", Version, "Anonymous")
	defaultBaseURL   = "https://oauth.reddit.com"
)

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:    httpClient,
		UserAgent: defaultUserAgent,
		BaseURL:   baseURL,
	}
	c.common.client = c
	c.Account = (*AccountService)(&c.common)
	c.Captcha = (*CaptchaService)(&c.common)
	c.Flair = (*FlairService)(&c.common)
	c.Gold = (*GoldService)(&c.common)
	c.LinksComments = (*LinksCommentsService)(&c.common)
	c.Listings = (*ListingsService)(&c.common)
	c.LiveThreads = (*LiveThreadsService)(&c.common)
	c.PrivateMessages = (*PrivateMessagesService)(&c.common)
	c.Misc = (*MiscService)(&c.common)
	c.Moderation = (*ModerationService)(&c.common)
	c.Multis = (*MultisService)(&c.common)
	c.Search = (*SearchService)(&c.common)
	c.Subreddits = (*SubredditsService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.Wiki = (*WikiService)(&c.common)
	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	reqURL := c.BaseURL.ResolveReference(rel)
	reqQuery := reqURL.Query()
	reqQuery.Set("raw_json", "1")
	reqURL.RawQuery = reqQuery.Encode()
	req, err := http.NewRequest(method, reqURL.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()
	if err := CheckResponse(resp); err != nil {
		return resp, err
	}
	if v != nil {
		return resp, json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, nil
}

func (c *Client) updateRateLimit(resp *http.Response) error {
	var err error
	t := time.Now()
	rl := rateLimit{}
	rl.Used, err = ratelimitGetInt(resp, "X-Ratelimit-Used")
	if err != nil {
		return err
	}
	rl.Remaining, err = ratelimitGetInt(resp, "X-Ratelimit-Remaining")
	if err != nil {
		return err
	}
	resetSec, err := ratelimitGetInt(resp, "X-Ratelimit-Reset")
	if err != nil {
		return err
	}
	rl.Reset = t.Add(time.Duration(resetSec) * time.Second)
	c.rateLimitMu.Lock()
	defer c.rateLimitMu.Unlock()
	c.rateLimit = &rl
	return nil
}

func (c *Client) RateLimitHit() bool {
	c.rateLimitMu.Lock()
	defer c.rateLimitMu.Unlock()
	if c.rateLimit == nil {
		return false
	}
	if c.rateLimit.Reset.Before(time.Now()) {
		c.rateLimit = nil
		return false
	}
	return false
}

type (
	ResponseError struct {
		Error string
		http.Response
	}
)

// Status code reddit uses to indicate that rate limit has been hit.
const statusCodeRateLimit = 429

// CheckResponse checks the response for correct status codes
// and keeps track of rate limits and returns nil if the response
// if ok. It returns a RateLimitError if the API's rate limiting
// blocked the response.
func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	if resp.StatusCode == statusCodeRateLimit {
		rle, err := rateLimitErrorFromResp(resp)
		if err != nil {
			return err
		}
		return rle
	}
	apiErr := &APIError{}
	rawJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.NewDecoder(bytes.NewReader(rawJSON)).Decode(apiErr)
	if err != nil {
		apiErr.Message = string(rawJSON)
		apiErr.ErrorCode = resp.StatusCode
		return apiErr
	}
	return apiErr
}

// UserAgent returns a string in reddit's desired format,
// intended to use as 'UserAgent' header for http requests.
func UserAgent(platform, appID, version, username string) string {
	return fmt.Sprintf("%s:%s:%s (by %s)",
		platform, appID, version, username)
}

// HTTPTransport implements the http.RoundTripper interface.
// It set the User-Agent to the given one for each request.
type HTTPTransport struct {
	UserAgent    string
	RoundTripper http.RoundTripper
}

func WrapHTTPTransport(userAgent string, roundTripper http.RoundTripper) *HTTPTransport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}
	return &HTTPTransport{
		UserAgent:    userAgent,
		RoundTripper: roundTripper,
	}
}

func (t *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	return t.RoundTripper.RoundTrip(req)
}

type (
	service struct {
		client *Client
	}

	AccountService         service
	CaptchaService         service
	FlairService           service
	GoldService            service
	LinksCommentsService   service
	LiveThreadsService     service
	PrivateMessagesService service
	ModerationService      service
	MultisService          service
	SearchService          service
	SubredditsService      service
	UsersService           service
	WikiService            service
)

type Response struct {
	*http.Response
}

// APIError implements the Error interface and is used to
// parse the JSON from API errors ins reponses.
type APIError struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.ErrorCode, e.Message)
}

type rateLimiter struct {
	err error
	sync.RWMutex
}
