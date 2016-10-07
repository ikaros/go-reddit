package reddit

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestResponseTypeAccount(t *testing.T) {
	rawJSON := bytes.NewBufferString(`{
		"kind": "t2",
		"data": {
			"has_mail": false,
			"name": "fooBar",
			"created": 123456789.0,
			"modhash": "f0f0f0f0f0f0f0f0...",
			"created_utc": 1315269998.0,
			"link_karma": 31,
			"comment_karma": 557,
			"is_gold": false,
			"is_mod": false,
			"has_verified_email": false,
			"id": "5sryd",
			"has_mod_mail": false
		}
	}`)

	var a Account
	if err := json.NewDecoder(rawJSON).Decode(&a); err != nil {
		t.Fatal(err)
	}
}
