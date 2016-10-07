package reddit

import (
	"fmt"
	"strings"
)

// ListingsService is the API Endpoint for listings
type ListingsService service

// ByID returns a listing of Links by fullname.
func (s *ListingsService) ByID(linkNames ...string) ([]Link, error) {
	for _, n := range linkNames {
		if !strings.HasPrefix(n, string(kindLink)) {
			return nil, fmt.Errorf("%s is no fullname of a link", n)
		}
	}
	r, err := s.client.NewRequest("GET", "/by_id/"+strings.Join(linkNames, ","), nil)
	if err != nil {
		panic(err)
	}
	var listing struct {
		Data struct {
			Children []struct {
				Data Link `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}
	if _, err := s.client.Do(r, &listing); err != nil {
		return nil, err
	}
	var links []Link
	for _, c := range listing.Data.Children {
		links = append(links, c.Data)
	}
	return links, err
}
