package reddit

import (
	"fmt"
	"net/http"
	"strconv"
)

// RateLimitError provides values for rate limiting.
type RateLimitError struct {
	// Approximate number of requests used in this period.
	// Header: X-Ratelimit-Used
	Used int

	// Approximate number of requests left to use.
	// Header: X-Ratelimit-Remaining
	Remaining int

	// Number of seconds to end of period.
	// Header: X-Ratelimit-Reset
	Reset int
}

func rateLimitErrorFromResp(resp *http.Response) (*RateLimitError, error) {
	var err error
	e := RateLimitError{}
	e.Used, err = ratelimitGetInt(resp, "X-Ratelimit-Used")
	if err != nil {
		return nil, err
	}
	e.Remaining, err = ratelimitGetInt(resp, "X-Ratelimit-Remaining")
	if err != nil {
		return nil, err
	}
	e.Reset, err = ratelimitGetInt(resp, "X-Ratelimit-Reset")
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func ratelimitGetInt(resp *http.Response, header string) (int, error) {
	h := resp.Header.Get(header)
	if h == "" {
		return 0, &RateLimitHeaderError{
			Response: resp,
			Message:  fmt.Sprintf("Header '%s' not found", header),
		}
	}
	i, err := strconv.Atoi(h)
	if err != nil {
		return 0, &RateLimitHeaderError{
			Response: resp,
			Message:  fmt.Sprintf("Header '%s': %s", header, err.Error()),
		}
	}
	return i, err
}

func (RateLimitError) Error() string {
	return "Rate limit has been hit"
}

// RateLimitHeaderError is returned by RateLimitFromResp to indicate
// that the given response include no, incorrect, or not necessary Headers.
type RateLimitHeaderError struct {
	*http.Response
	Message string
}

func (e *RateLimitHeaderError) Error() string {
	return e.Message
}
