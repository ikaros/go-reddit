package reddit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestListingsByIDAPIInternalError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "Server fuckup","error": 503}`)
	}))
	defer ts.Close()
	client := NewClient(nil)
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	client.BaseURL = u
	links, err := client.Listings.ByID("t3_asdfgh")
	if err == nil {
		t.Error("No error returned")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatal("No APIError returned")
	}

	if apiErr.Message != "Server fuckup" {
		t.Error("Wrong Message")
	}
	if apiErr.ErrorCode != 503 {
		t.Error("Wrong ErrorCode")
	}

	if len(links) > 0 {
		t.Error("Links and error returned")
	}
}
