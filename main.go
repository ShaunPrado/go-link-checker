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

func fetchAndParse(ud UrlDepth, ch chan<- UrlDepth, wg *sync.WaitGroup) {
	defer wg.Done()

	if ud.Depth > 3 {
		return
	}

	resp, err := http.Get(ud.Url)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

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
						absUrl := resp.Request.URL.ResolveReference(parsedUrl).String()
						ch <- UrlDepth{Url: absUrl, Depth: ud.Depth + 1}
						break
					}
				}
			}
		}
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
