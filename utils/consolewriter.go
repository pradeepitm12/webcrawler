package utils

import "fmt"

type ConsoleWriter struct{}
/**
* Write
* Simply writes data on console
 */
func (c ConsoleWriter) Write(data string) {
	fmt.Println(data)
}
