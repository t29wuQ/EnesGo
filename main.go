package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"./nes"
)

func main(){
	file, err := os.Open("./helloworld.nes")
	if err != nil{
		log.Println("Faild : load rom file")
		return
	}
	defer file.Close()
	rom, err := ioutil.ReadAll(file)
	if err != nil{
		log.Println("Faild : read rom file")
		return
	}
	fmt.Println(rom)

	//cpu := nes.CreateCpu(rom[0x10:0x10+int(rom[4])*0x4000])
	ppu := nes.CreatePpu(rom[int(0x10+int(rom[4])*0x4000):])
	ppu.Debug()
	
}