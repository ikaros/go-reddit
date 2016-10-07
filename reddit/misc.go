package reddit

import "encoding/json"

// MiscService is the API Endpoint for misc things
type MiscService service

// Scopes returns available OAuth2 scopes.
// If no scopes are given, information on all scopes are returned.
// Invalid scope(s) will result in a 400 error with body that
// indicates the invalid scope(s).
func (s *MiscService) Scopes(scope string) (*Scopes, *Response, error) {
	r, err := s.client.NewRequest("GET", "/api/v1/scopes", nil)
	if err != nil {
		panic(err)
	}
	if scope != "" {
		r.URL.Query().Add("scope", scope)
	}
	resp, err := s.client.client.Do(r)
	if err != nil {
		return nil, nil, err
	}
	rp := &Response{Response: resp}
	scopes := make(Scopes)
	if err := json.NewDecoder(resp.Body).Decode(&scopes); err != nil {
		return nil, rp, err
	}
	return &scopes, rp, nil
}

// OAuth permissions scopes for authentication.
type Scopes map[string]struct {
	Description string
	ID          string
	Name        string
}
