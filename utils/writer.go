package utils

import (
	"fmt"
	"os"
	)

type Writer interface{
	Write(string)
}

type ConsoleWriter struct {}

func(c ConsoleWriter)Write(data string){
	fmt.Println(data)
}

type FileWriter struct {
	FilePath string
	//FilePtr *os.File

}
//func(f FileWriter)Init(){
//	file,err:=os.OpenFile(f.FilePath, os.O_RDWR|os.O_APPEND, 0660)
//	//file,err:=os.Create(f.FilePath)
//	if err!=nil{
//		fmt.Println("Cannot Create file",err)
//	}
//	f.FilePtr=file
//}
func(f FileWriter)Write(data string){
	file, err := os.Create(f.FilePath)
	defer file.Close()
	_,err=file.WriteString(data)
	fmt.Println("Message while writing ", data)
	if err!=nil{
		fmt.Println("Error in writing " ,err)
	}
	file.Sync()
}