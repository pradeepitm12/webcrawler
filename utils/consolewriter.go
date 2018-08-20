package utils

import "fmt"

type ConsoleWriter struct{}

/**
* Write
* Simply writes data on console
 */
func (c ConsoleWriter) Write(datapipe <-chan string) {
	for data := range datapipe {
		fmt.Println(data)
	}
}
