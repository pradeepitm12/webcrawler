package main

import (
	"flag"
	"fmt"
	"os"

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
	redisWriter := utils.RedisWriter{RedisAddress: "127.0.0.1:6379", RedisPassword: "redis@123@Azure", RedisPoolIdleTimeout: 100, RedisPoolMaxActive: 100, RedisPoolMaxIdle: 100}
	redisWriter.InitRedisPool()
	var w utils.Writer = redisWriter
	// initilize a buffered channel of string
	input := make(chan string, 10)
	output := make(chan []string, 10)
	writerPipe := make(chan string, 10)
	const numberOfWorkers = 3
	for i := 0; i < numberOfWorkers; i++ {
		go utils.CrawlWorker(input, output)
		//for j := 0; j < numberOfWorkers; j++ {
		go w.Write(writerPipe)
		//}
	}

	for _, link := range rootLink {
		input <- link
	}
	for linkarray := range output {
		for _, link := range linkarray {
			fmt.Println("ready to be written --- ", link)
			writerPipe <- link
		}
	}
}
