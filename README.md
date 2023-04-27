# Go Site Crawler

[![CI/CD](https://github.com/ShaunPrado/go-link-checker/actions/workflows/go.yml/badge.svg)](https://github.com/ShaunPrado/go-link-checker/actions/workflows/go.yml)


A simple, multithreaded site crawler and link checker in Go that crawls a given website and prints the visited URLs. The crawler operates up to a maximum depth of 3 levels and can be configured to stop after visiting a specified number of URLs.

## Features

* Multithreaded crawling with configurable worker limits
* Tracks visited URLs to avoid redundant crawling
* Stops crawling after reaching a specified depth or visiting a certain number of URLs

Usage
Clone this repository:

```
git clone https://github.com/ShaunPrado/go-link-checker.git
cd go-link-checker
```
Build the project:

```
go build -o main
```
Run the crawler with the required flags:


```
./main -site="https://example.com" -totalUrls=100
```
Replace https://example.com with the website you want to crawl, and 100 with the total number of URLs you want the crawler to visit before stopping.

## Command Line Flags
- site: The website URL to start crawling from (e.g., "https://example.com"). This flag is required.
- totalUrls: The total number of URLs the crawler should visit before stopping. Default: 100.
## Dependencies
golang.org/x/net/html: For parsing HTML content.

## Roadmap

- Additional browser support

- Add more integrations

## License
This project is open-source and available under the MIT License.


# Short-term goals

* Parallel crawling: Enhance the crawler's performance by introducing parallel crawling for multiple websites or sections of the same site.

# Long-term goals

* Distributed crawling architecture to scale the crawling process across multiple cloud instances.
* Web-based user interface for easier configuration and monitoring of the crawler.
* Allow users to schedule crawls to run at specific intervals.