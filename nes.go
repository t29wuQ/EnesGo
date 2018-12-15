package main

import(
	"fmt"
	"syscall/js"
	"strings"
	"strconv"
	"os"
	"./nes"
	"encoding/base64"
	"bytes"
	"image"
	"image/jpeg"
	"image/color"
)

func convertStringArrayToBinaryArray(array []string) []uint8{
	binary := make([]uint8, len(array))
	for i := 0;i < len(array);i++{
		num, _ := strconv.Atoi(array[i])
		binary[i] = uint8(num)
	}
	return binary
}

func main(){
	fmt.Println("Hello")
	document := js.Global().Get("document")
	binaryString := document.Call("getElementById", "binary").Get("innerText").String()
	binary := strings.Split(binaryString, ",")
	fmt.Println(binary[1])
	if binary[0] != "78" || binary[1] != "69" || binary[2] != "83" || binary[3] != "26"{
		fmt.Println("error: cant load nes rom")
		os.Exit(0)
	}
	header := convertStringArrayToBinaryArray(binary[:0x10])
	ppu := nes.CreatePpu(convertStringArrayToBinaryArray(binary[int(0x10+int(header[4])*0x4000):]))
	rawImage := ppu.Debug()
	img := image.NewNRGBA(image.Rect(0,0,8,8))
	for i := 0;i < 8;i++{
		for j := 0;j < 8;j++{
			k := i * 24 + j * 3
			img.Set(j, i, color.RGBA{rawImage[k], rawImage[k + 1], rawImage[k + 2], 255})
		}
	}
	var buffer bytes.Buffer
	jpeg.Encode(&buffer, img, nil)
	imageEnc := base64.StdEncoding.EncodeToString(buffer.Bytes())
	jsimg := document.Call("getElementById", "screen")
	jsimg.Set("src", "data:image/jpeg;base64,"+imageEnc)
}