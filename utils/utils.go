package utils

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type PipeLine struct {
	Input     chan string
	Output    chan []string
	Write     chan string
	CheckHash map[string]int
}

func GetPipeLine() *PipeLine {
	return &PipeLine{make(chan string, 10), make(chan []string, 10), make(chan string, 5), make(map[string]int)}
}

/**
* GetLinks
* For the body received
* if token of body is starttoken and is of <a>
* itterate over all the attributes and look for href to get the links
* put all the links in map with value as 1 :- this ensures that links are not repeted
* return all the links as an array
 */
func (pipe *PipeLine) GetLinks(data io.Reader) []string {
	links := []string{}
	//linkHash := pipe.checkHash
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
					aV := attribute.Val
					if pipe.CheckHash[aV] == 0 && !strings.HasPrefix(aV, "#") && !strings.HasPrefix(aV, "/") {
						pipe.CheckHash[aV] = 1
						links = append(links, aV)
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

func (pipe *PipeLine) Work(link string) []string {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}

	resp, err := client.Get(CheckLink(link))
	if err != nil {
		log.Println("Error while client.Get(link) : ", err)
	}
	defer resp.Body.Close()
	return pipe.GetLinks(resp.Body)
}

/**
* CrawlWorker
* takes job as string pass it to work
* string array from work is passed to string array channel
 */
func (pipe *PipeLine) CrawlWorker() {
	for work := range pipe.Input {
		pipe.Output <- pipe.Work(work)
	}
}

func CheckLink(link string) string {
	return strings.Split(link, ";")[0]
}
