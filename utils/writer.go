package utils

/**
* writer
* Writer
* A behaviour which takes a string input channel
 */
type Writer interface {
	Write(<-chan string)
}
