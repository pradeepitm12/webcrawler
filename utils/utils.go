package utils

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func GetLinks(data io.Reader) []string {
	links := []string{}
	//temp:= []string{}
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
					if linkHash[attribute.Val] == 0 {
						linkHash[attribute.Val] = 1
						links = append(links, attribute.Val)
					}
				}
			}
		}
	}
}

func LinkReader(link string) []string {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(link)
	if err != nil {
		//fmt.Println("Error while client.Get(link) : ", err)
		log.Fatal("Error while client.Get(link) : ", err)
	}
	defer resp.Body.Close()
	return GetLinks(resp.Body)

}
