package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResponseBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Test response")
	}))
	defer ts.Close()

	body, err := getResponseBody(ts.URL)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Test response")
}

func TestParseLinks(t *testing.T) {
	baseUrl, _ := url.Parse("http://example.com")
	testHTML := `
		<html>
		<body>
			<a href="/test1">Test 1</a>
			<a href="/test2">Test 2</a>
			<a href="/test3">Test 3</a>
		</body>
		</html>
	`

	links, err := parseLinks([]byte(testHTML), baseUrl)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(links))
	assert.Contains(t, links, "http://example.com/test1")
	assert.Contains(t, links, "http://example.com/test2")
	assert.Contains(t, links, "http://example.com/test3")
}

func TestFetchAndParse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
			<html>
			<body>
				<a href="/test1">Test 1</a>
				<a href="/test2">Test 2</a>
				<a href="/test3">Test 3</a>
			</body>
			</html>
		`)
	}))
	defer ts.Close()

	ch := make(chan UrlDepth)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		fetchAndParse(UrlDepth{Url: ts.URL, Depth: 1}, ch, &wg)
		close(ch)
	}()

	var count int
	for ud := range ch {
		assert.True(t, strings.HasPrefix(ud.Url, ts.URL))
		count++
	}
	assert.Equal(t, 3, count)
}
