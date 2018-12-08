package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
)

func main(){
	file, err := os.Open("./helloworld.nes")
	if err != nil{
		log.Println("Faild : load rom file")
		return
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil{
		log.Println("Faild : read rom file")
		return
	}
	fmt.Println(buf)
	fmt.Println("hello world")
}