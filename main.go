package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func fetchAndParse(url string, depth int, ch chan<- [2]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth > 3 {
		return
	}

	resp, err := http.Get(url)
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
						ch <- [2]interface{}{a.Val, depth + 1}
						break
					}
				}
			}
		}
	}
}

func main() {
	sitePtr := flag.String("site", "default", "Input string")
	totalUrlsPtr := flag.Int("totalUrls", 100, "Total number of URLs to visit")
	flag.Parse()

	// Get the value of the "input" flag
	startURL := *sitePtr

	// Print the input value
	fmt.Println("Input string:", startURL)

	ch := make(chan [2]interface{})
	var wg sync.WaitGroup

	// Limit the number of concurrent workers
	workerLimit := make(chan struct{}, 10)

	wg.Add(1)
	go func() {
		workerLimit <- struct{}{}
		fetchAndParse(startURL, 1, ch, &wg)
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

		url := data[0].(string)
		currentDepth := data[1].(int)
		if !visited[url] {
			visited[url] = true
			fmt.Println(url)
			urlsVisited++
			wg.Add(1)
			go func(url string, depth int) {
				workerLimit <- struct{}{}
				fetchAndParse(url, depth, ch, &wg)
				<-workerLimit
			}(url, currentDepth)
		}
	}
	// Print the visited URLs
	fmt.Println("length:", len(visited))

}
