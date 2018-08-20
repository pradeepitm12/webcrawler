package utils

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

/**
* GetLinks
* For the body received
* if token of body is starttoken and is of <a>
* itterate over all the attributes and look for href to get the links
* put all the links in map with value as 1 :- this ensures that links are not repeted
* return all the links as an array
 */
func GetLinks(data io.Reader) []string {
	links := []string{}
	linkHash := make(map[string]int)
	page := html.NewTokenizer(data)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attribute := range token.Attr {
				if attribute.Key == "href" {
					//add link only if not present
					if linkHash[attribute.Val] == 0 && !strings.HasPrefix(attribute.Val, "#") && !strings.HasPrefix(attribute.Val, "/") {
						linkHash[attribute.Val] = 1
						links = append(links, attribute.Val)
					}
				}
			}
		}
	}
}

/**
* LinkReader
* do a get request on the link
* pass the request body to get all the links as array
 */

func Work(link string) []string {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(link)
	if err != nil {
		log.Fatal("Error while client.Get(link) : ", err)
	}
	defer resp.Body.Close()
 	return  GetLinks(resp.Body)
}
/**
* CrawlWorker
* takes job as string pass it to work
* string array from work is passed to string array channel
 */
func CrawlWorker(input <-chan string, output chan<- []string) {
	for work := range input {
		output <- Work(work)
	}
}
