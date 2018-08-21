package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/pradeepitm12/webcrawler/utils"
)

//var wg = sync.WaitGroup{}

/**
* for each root links
* collect all the links as an array of strings
* Instanciate  a writer type
* *****Example to write in File
* fileWriter:=utils.FileWriter{FilePath:"/home/zaid/Desktop/Test.txt"}
* *****Example to write on Console
* var w utils.Writer = utils.ConsoleWriter{}
* *****Example to write in redis
* redisWriter:=utils.RedisWriter{RedisAddress: "127.0.0.1:6379", RedisPassword: "redis@123",
* RedisPoolIdleTimeout:100,RedisPoolMaxActive:100,RedisPoolMaxIdle:100}
* redisWriter.InitRedisPool()
 */
func main() {
	fmt.Println(" ****** This is a web crawler ******")
	flag.Parse()
	//take console input space seperated root links
	rootLink := flag.Args()
	//check if no input if given
	if len(rootLink) < 1 {
		fmt.Println("Please provide a link")
		os.Exit(1)
	}
	/**
	* writing links to redis,
	* links can be written in file or console as well.
	 */

	var w utils.Writer = utils.FileWriter{FilePath: "/home/pradheep/Desktop/Crawler.txt"}

	//pipeLine := utils.PipeLine{make(chan string, 10), make(chan []string, 10), make(chan string, 5)}
	pipeLine := utils.GetPipeLine()
	//fmt.Println("First ", pipeLine)
	const numberOfWorkers = 3
	for i := 0; i < numberOfWorkers; i++ {
		go pipeLine.CrawlWorker()
		for j := 0; j < numberOfWorkers; j++ {
			go w.Write(pipeLine.Write)
		}
	}

	for _, link := range rootLink {
		defer close(pipeLine.Input)
		pipeLine.Input <- link
	}
	for linkarray := range pipeLine.Output {
		defer close(pipeLine.Output)
		for _, link := range linkarray {
			defer close(pipeLine.Write)
			// adding check for host name
			// start putting new links to input if they have redhat has hostname
			rawUrl, _ := url.Parse(link)
			if strings.Contains(rawUrl.Hostname(), "red") {
				fmt.Println("Writing link ", link)
				pipeLine.Input <- link
				//	fmt.Println("added to input   ", link)
			}
			pipeLine.Write <- link
		}
	}
	//fmt.Println("First ", pipeLine)
}
