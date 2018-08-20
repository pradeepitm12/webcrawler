package utils

/**
* writer
* Write
* A behaviour which takes a string input channel
 */
type Writer interface {
	Write(<-chan string)
}
