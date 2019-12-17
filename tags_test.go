package cloudscale

import (
	"fmt"
	"net/http"
	"testing"
)

var toQueryStringTestCases = []struct {
	tags     TagMap
	expected string
}{
	{TagMap{"a": "b"}, "http://example.com?tag%3Aa=b"},
	{TagMap{"a": ""}, "http://example.com?tag%3Aa="},
	{TagMap{"a": "b", "c": "d"}, "http://example.com?tag%3Aa=b&tag%3Ac=d"},
}
func TestTagsToQueryString(t *testing.T) {
	for _, tt := range toQueryStringTestCases {
		t.Run(fmt.Sprintf("%#v", tt.tags), func(t *testing.T) {
			// arrange
			req, _ := http.NewRequest("GET", "http://example.com", nil)

			// act
			requestModifier := WithTagFilter(tt.tags)
			requestModifier(req)

			// assert
			if acutal := req.URL.String(); acutal != tt.expected {
				t.Errorf("Unexpected result\n got=%#v\nwant=%#v", acutal, tt.expected)
			}
		})
	}
}
