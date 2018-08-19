package utils

import (
	"os"
	"log"
)

type FileWriter struct {
	FilePath string
}

/**
* Write
* creates a file is not created, and keep appending data to new line
 */
func (f FileWriter) Write(data string) {
	file,_:=os.OpenFile(f.FilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer file.Close()
	_, err := file.WriteString(data+"\n")
	log.Println("Message while writing ", data)
	if err != nil {
		log.Println("Error in writing ", err)
	}
	file.Sync()
}