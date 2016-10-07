package reddit

import (
	"net/http"
	"testing"
)

func TestRateLimitFetchFromRespError(t *testing.T) {
	resp := &http.Response{}
	if _, err := rateLimitErrorFromResp(resp); err == nil {
		t.Error("No error has been returned")
	}
}

func TestRateLimitFetchFromResp(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"X-Ratelimit-Used":      {"10"},
			"X-Ratelimit-Remaining": {"10"},
			"X-Ratelimit-Reset":     {"10"},
		},
	}

	rl, err := rateLimitErrorFromResp(resp)
	if err != nil {
		t.Fatal(err)
	}
	if i := rl.Used; i != 10 {
		t.Errorf("Used was %d instead of 10", i)
	}
	if i := rl.Remaining; i != 10 {
		t.Errorf("Remaining was %d instead of 10", i)
	}
	if i := rl.Reset; i != 10 {
		t.Errorf("Rest was %v instead of 10", i)
	}
}

func TestRateLimitFetchFromConversion(t *testing.T) {
	for i, header := range []http.Header{
		{
			"X-Ratelimit-Used":      {"Foo"},
			"X-Ratelimit-Remaining": {"10"},
			"X-Ratelimit-Reset":     {"10"},
		},
		{
			"X-Ratelimit-Used":      {"10"},
			"X-Ratelimit-Remaining": {"Foo"},
			"X-Ratelimit-Reset":     {"10"},
		},
		{
			"X-Ratelimit-Used":      {"10"},
			"X-Ratelimit-Remaining": {"10"},
			"X-Ratelimit-Reset":     {"Foo"},
		},
		{
			"X-Ratelimit-Used":      {"10.0"},
			"X-Ratelimit-Remaining": {"10"},
			"X-Ratelimit-Reset":     {"10"},
		},
		{
			"X-Ratelimit-Used":      {"10"},
			"X-Ratelimit-Remaining": {"10.0"},
			"X-Ratelimit-Reset":     {"10"},
		},
		{
			"X-Ratelimit-Used":      {"10"},
			"X-Ratelimit-Remaining": {"10"},
			"X-Ratelimit-Reset":     {"10.0"},
		},
	} {
		_, err := rateLimitErrorFromResp(&http.Response{Header: header})
		if err == nil {
			t.Errorf("Test(%d) Returned no error for broken header(%#v)",
				i, header)
		}
	}
}
