package main

import (
	"fmt"
	"flag"
	"github.com/pradeepitm12/webcrawler/utils"
	"os"
	"sync"
)
var wg = sync.WaitGroup{}

func main(){
	fmt.Println(" ****** This is a web crawler ******")
	flag.Parse()
	//take console input space seperated root links
	rootLink := flag.Args()
	//check if no input if given
	if len(rootLink) <1{
		fmt.Println("Please provide a link")
		os.Exit(1)
	}
	//for each root links
	//collect all the links as an array of strings
	//Instanciate  a writer type

	//fileWriter:=utils.FileWriter{FilePath:"/home/zaid/Desktop/Test.txt"}
	//fileWriter.Init()
	var w utils.Writer = utils.ConsoleWriter{}


	pipe:=make(chan string,10)
	for _,link:=range rootLink{
		/**
		* 1. Write on console.
		* 2. Write in File
		* X. For future use like Write in Redis, Kafka, or what ever is required.
		*/
		pipe<-link
		wg.Add(1)
		go  DataTransfer(utils.LinkReader(<-pipe),w)

	}
	wg.Wait()
}
func DataTransfer(data []string, writer utils.Writer) {

	for _,link:=range data{
		writer.Write(link)
	}
	wg.Done()
}