package gcache

import (
	"testing"
)

var vtests = []struct {
	url      string
	filename string
}{
	{"http://example.com/",
		"http://example.com/"},
	{"https://example.com/index.html",
		"https://example.com/index.html"},
	{"http://example.com:80/",
		"http://example.com/"},
	{"https://example.com:443/",
		"https://example.com/"},
	{"http://example.com:8080/",
		"http://example.com:8080/"},
	{"http://example.com:443/",
		"http://example.com:443/"},
	{"https://example.com:80/",
		"https://example.com:80/"},
}

func TestGetFileName(t *testing.T) {
	for _, vt := range vtests {
		if filename, err := GetFileName(vt.url); err != nil {
			t.Errorf("%q, got %q", vt.url, err)
		} else if filename != vt.filename {
			t.Errorf("%q, want %q", filename, vt.filename)
		}
	}
}
