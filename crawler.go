package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/pradeepitm12/webcrawler/utils"
	)

var wg = sync.WaitGroup{}
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
	redisWriter:=utils.RedisWriter{RedisAddress: "127.0.0.1:6379", RedisPassword: "redis@123",RedisPoolIdleTimeout:100,RedisPoolMaxActive:100,RedisPoolMaxIdle:100}
	redisWriter.InitRedisPool()
	var w utils.Writer =redisWriter
	// initilize a buffered channel of string
	pipe := make(chan string, 10)
	// iterate through all the links given by command line
	for _, link := range rootLink {
		//putting links in the pipe
		pipe <- link
		// adding the number of process in wait group so main knows when to stop
		wg.Add(1)
		//go routine started for items in pipe
		go DataTransfer(utils.LinkReader(<-pipe), w)

	}
	// main waits till last item in wait group
	wg.Wait()
}
/**
* DataTransfer
* For slice of an array of links
* write all the links in the type of writer received
 */
func DataTransfer(data []string, writer utils.Writer) {

	for _, link := range data {
		writer.Write(link)
	}
	wg.Done()
}
