package nes

import(
	"fmt"
	"strconv"
)

type Ppu struct{
	sprite [512][8][8]uint8
	vram [0x4000]uint8
}

func CreatePpu(crom []uint8) *Ppu{
	ppu := new(Ppu)
	loopCount := 0
	for _, b := range crom {
		ppu.vram[loopCount] = b
		loopCount++;
	}
	for i := 0;i < len(crom)/16; i++{
		for j := 0; j < 8;j++{
			low := []rune(formatBynary(int(crom[i * 16 + j])))
			high := []rune(formatBynary(int(crom[i * 16 + j + 8])))
			for k := 0;k < 8;k++{
				lbit, _ := strconv.Atoi(string(low[k]))
				hbit, _ := strconv.Atoi(string(high[k]))
				ppu.sprite[i][j][k] = uint8(hbit * 2 + lbit)
			}
		}
	}
	return ppu
}

func (ppu Ppu) Debug(){
	n := 73
	for i := 0;i < 8;i++{
		for j := 0;j < 8;j++{
			if ppu.sprite[n][i][j] > 0{
				fmt.Print("■")
			}else{
				fmt.Print("□")
			}
		}
		fmt.Println("")
	}
}

func formatBynary(i int) string{
	str := strconv.FormatInt(int64(i), 2)
	for j := len(str);j < 8;j++{
		str = "0" + str
	}
	return str
}