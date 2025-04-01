package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

func isValidURL(str string) bool {
	parsedURL, err := url.ParseRequestURI(str)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

func handleCommands(input string) {
	commandTokens := strings.Fields(input)
	if len(commandTokens) == 0 {
		fmt.Println("Unknown command. Run 'scrapper help' for usage.")
		return
	}
	if commandTokens[0] == "help" {
		fmt.Println("Commands:")
		fmt.Println("  exit - Exits the program")
		fmt.Println("  scrapper [url of the website to be scrapped] - recursively navigates the website to check for deadlinks")
	}
	if commandTokens[0] == "scrapper" {
		if len(commandTokens) > 1 {
			if isValidURL(commandTokens[1]) {
				scrape(commandTokens[1])
				return
			} else {
				fmt.Println(commandTokens[1], "is NOT a valid URL ❌")
				return
			}
		}
	}

	fmt.Println("Unknown command. Run 'help' for usage.")
}

func scrape(url string) {
	type urlStatus struct {
		url    string
		status string
	}
	var urlsChecked []urlStatus
	resp, err := http.Get(url)
	if err != nil {
		urlsChecked = append(urlsChecked, urlStatus{url, err.Error()})
		return
	}
	defer resp.Body.Close()

	// Check if status is 200 OK
	if resp.StatusCode != http.StatusOK {
		urlsChecked = append(urlsChecked, urlStatus{url, resp.Status})
		return
	}

	// Check if Content-Type is HTML
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		urlsChecked = append(urlsChecked, urlStatus{url, "Non HTML content type"})
		return
	}
	urlsChecked = append(urlsChecked, urlStatus{url, resp.Status})

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var links []string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// Find the href attribute
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		// Recursively traverse child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)
	fmt.Println(links)
}
