package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	var (
		webPageUrl string
		links      []string
	)
	fmt.Println("Enter valid html webPageUrl:")
	fmt.Fscan(os.Stdin, &webPageUrl)
	resp, err := parseFromSite(webPageUrl)
	if err != nil {
		log.Fatalf("error fetching URL: %v\n", err)
	}
	links, err = getAllLinks(resp.Body, webPageUrl)
	fmt.Println("Parsed Links:")
	for _, val := range links {
		fmt.Println(val)
	}
}
func parseFromSite(websiteUrl string) (*http.Response, error) {
	resp, err := http.Get(websiteUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code was %d\n", resp.StatusCode)
	}
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}
	return resp, err
}

func getAllLinks(body io.Reader, siteUrl string) ([]string, error) {
	var links []string
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" && strings.HasPrefix(attr.Val, "/") {
						urlPage, err := url.Parse(siteUrl)
						if err != nil {
							return nil, err
						}
						finalUrl, err := urlPage.Parse(attr.Val)
						links = append(links, finalUrl.String())
					} else if attr.Key == "href" && !strings.HasPrefix(attr.Val, "#") {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
	return links, nil
}
