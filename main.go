package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	fmt.Println("Welcome to the GoLang Deadlink-Scrapper v1.")
	fmt.Println("For help run 'help'")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter commands (type 'exit' to quit):")

	for {
		fmt.Print("> ") // Prompt for input
		scanner.Scan()  // Read user input
		input := strings.TrimSpace(scanner.Text())
		words := strings.Fields(input)
		switch words[0] {
		case "exit":
			fmt.Println("Exiting program...")
			break // Stop the loop
		case "help":
			fmt.Println("Usage:")
			fmt.Println("scrape <url>")
		case "scrape":
			if len(words) > 1 {
				handleUrl(words[1])
			} else {
				fmt.Println("Unknown command")
			}

		}
	}

	fmt.Println("")
}

type Crawler struct {
	visited map[string]bool
	mu      sync.Mutex
	wg      sync.WaitGroup
	results chan string
}

// Fetches a page, extracts internal links, and recursively crawls
func (c *Crawler) crawl(targetURL string) {
	defer c.wg.Done()

	// Mark URL as visited
	c.mu.Lock()
	if c.visited[targetURL] {
		c.mu.Unlock()
		return
	}
	c.visited[targetURL] = true
	c.mu.Unlock()

	// Check link status
	resp, err := http.Get(targetURL)
	if err != nil {
		c.results <- fmt.Sprintf("%s - ERROR: %s", targetURL, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		c.results <- fmt.Sprintf("%s - %d %s", targetURL, resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	// Only parse HTML pages
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		c.results <- fmt.Sprintf("%s - %s", targetURL, "", "Is not html")
		return
	}
	c.results <- fmt.Sprintf("%s - %d %s", targetURL, resp.StatusCode, http.StatusText(resp.StatusCode))

	// Parse links from HTML
	links, err := extractInternalLinks(targetURL, resp)
	if err != nil {
		return
	}

	// Crawl discovered links
	for _, link := range links {
		c.wg.Add(1)
		go c.crawl(link)
	}
}

// Extracts internal links from an HTML document
func extractInternalLinks(baseURL string, resp *http.Response) ([]string, error) {
	var links []string
	parsedBase, _ := url.Parse(baseURL)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := attr.Val
					parsedLink, err := url.Parse(link)
					if err == nil && parsedLink.Host == "" {
						absoluteURL := parsedBase.ResolveReference(parsedLink).String()
						links = append(links, absoluteURL)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)

	return links, nil
}

// Checks if the response is an HTML page
func isHTML(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	return contentType == "text/html" || contentType == "text/html; charset=UTF-8"
}

func handleUrl(startUrl string) {
	//startUrl := "https://scrape-me.dreamsofcode.io" // Change this to the target website

	crawler := &Crawler{
		visited: make(map[string]bool),
		results: make(chan string, 100),
	}

	fmt.Println("Starting crawl...")

	// Start crawling
	crawler.wg.Add(1)
	go crawler.crawl(startUrl)

	// Goroutine to close results channel when done
	go func() {
		crawler.wg.Wait()
		close(crawler.results)
	}()

	// Print results as they arrive
	for result := range crawler.results {
		fmt.Println(result)
	}

	fmt.Println("Crawl complete.")
}
