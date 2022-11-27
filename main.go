package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	var (
		webPageUrl string
	)
	fmt.Println("Enter valid html webPageUrl:")
	fmt.Fscan(os.Stdin, &webPageUrl)
	resp := parseFromSite(webPageUrl)
	compiled, err := regexp.Compile("<a.*?href=[\"\"'](?P<url>[^\"\"^']+[.]*?)[\"\"'].*?>(?P<keywords>[^<]+[.]*?)</a>")
	if err != nil {
		log.Fatalf("%s", err)
	}
	data, _ := io.ReadAll(resp.Body)
	all := compiled.FindAllStringSubmatch(string(data), -1)
	fmt.Print("\n-------------------------------PARSED LINKS----------------------------------")
	for _, slice1 := range all {
		fmt.Printf("\n")
		if !strings.HasPrefix(slice1[1], "/") {
			urlPage, err := url.Parse(webPageUrl)
			if err != nil {
				log.Fatalf("Cant parse main page of webUrl")
			}
			finalUrl, err := urlPage.Parse(slice1[1])
			fmt.Print(finalUrl)
		} else {
			fmt.Print(slice1[1])
		}
	}
	fmt.Println("\n-----------------------------------------------------------------------------")
}

func parseFromSite(url string) *http.Response {
	htmlDoc, err := http.Get(url)
	if err != nil {
		log.Fatalf("Invalid URL: %s", err)
	}
	if htmlDoc.StatusCode != http.StatusOK {
		log.Fatalf("Cant connect to site, err code: %s", err)
	}
	ctype := htmlDoc.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}
	return htmlDoc
}
