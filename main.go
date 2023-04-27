package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"
)

type UrlDepth struct {
	Url   string
	Depth int
}

func getResponseBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func parseLinks(body []byte, base *url.URL) ([]string, error) {
	var links []string
	z := html.NewTokenizer(bytes.NewReader(body))

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						parsedUrl, err := url.Parse(a.Val)
						if err != nil {
							continue
						}
						absUrl := base.ResolveReference(parsedUrl).String()
						links = append(links, absUrl)
						break
					}
				}
			}
		}
	}

	return links, nil
}

func fetchAndParse(ud UrlDepth, ch chan<- UrlDepth, wg *sync.WaitGroup) {
	defer wg.Done()

	if ud.Depth > 3 {
		return
	}

	body, err := getResponseBody(ud.Url)
	if err != nil {
		fmt.Println(err)
		return
	}

	respUrl, err := url.Parse(ud.Url)
	if err != nil {
		fmt.Println(err)
		return
	}

	links, err := parseLinks(body, respUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, link := range links {
		ch <- UrlDepth{Url: link, Depth: ud.Depth + 1}
	}
}

func main() {
	sitePtr := flag.String("site", "", "Input string")
	totalUrlsPtr := flag.Int("totalUrls", 100, "Total number of URLs to visit")
	flag.Parse()

	if *sitePtr == "" {
		fmt.Println("Please provide a valid site URL using the -site flag.")
		return
	}

	ch := make(chan UrlDepth)
	var wg sync.WaitGroup
	workerLimit := make(chan struct{}, 10)

	wg.Add(1)
	go func() {
		workerLimit <- struct{}{}
		fetchAndParse(UrlDepth{Url: *sitePtr, Depth: 1}, ch, &wg)
		<-workerLimit
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Keep track of visited URLs to avoid redundant crawling
	visited := make(map[string]bool)
	urlsVisited := 0

	for data := range ch {
		if urlsVisited >= *totalUrlsPtr {
			break
		}

		if !visited[data.Url] {
			visited[data.Url] = true
			urlsVisited++
			wg.Add(1)
			go func(ud UrlDepth) {
				workerLimit <- struct{}{}
				fetchAndParse(ud, ch, &wg)
				<-workerLimit
			}(data)
		}
	}
	fmt.Println("Visited URLs count:", len(visited))
}
