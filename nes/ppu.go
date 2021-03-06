package nes

import(
	"strconv"
)

func formatBynary(i int) string{
	str := strconv.FormatInt(int64(i), 2)
	for j := len(str);j < 8;j++{
		str = "0" + str
	}
	return str
}

type Ppu struct{
	sprite [512][8][8]uint8
	vram [0x4000]uint8
	nmiInterrupt bool
	spriteSize uint8
	bgPatternTable int
	spritePatternTable int
	displayTable int
	ppuAddressInc uint16
	isSpritevVisible bool
	isBgVisible bool
	oamAddress uint8
	oamWriteCount int
	tmpOam [4]uint8
	writeCount int
	scrollOffsetX uint8
	scrollOffsetY uint8
	ppuAddress uint16
	vBlank bool
	spriteHit bool

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

func (ppu Ppu) Debug() []uint8{
	image := make([]uint8, 192)
	//image := ""
	n := 73
	count := 0
	for i := 0;i < 8;i++{
		for j := 0;j < 8;j++{
			if ppu.sprite[n][i][j] > 0{
				image[count] = 0xff
				image[count + 1] = 0xff
				image[count + 2] = 0xff
				//image += "255,255,255,"
			}else{
				image[count] = 0
				image[count + 1] = 0
				image[count + 2] = 0
				//image += "0,0,0,"
			}
			count += 3
		}
	}
	return image
}

func (ppu Ppu) readVram(address uint16) uint8{
	return ppu.vram[address]
}

func (ppu Ppu) writePpuRegister(address uint16, value uint8){
	switch address{
	case 0x2000:
		ppu.nmiInterrupt = refbit(value, 7) == 1
		ppu.spriteSize = refbit(value, 5) * 8 + 8
		ppu.bgPatternTable = 256 * int(refbit(value, 4))
		ppu.spritePatternTable = 256 * int(refbit(value, 3))
		ppu.displayTable = int(refbit(value, 1)) * 2 + int(refbit(value, 0))
	case 0x2001:
		ppu.isSpritevVisible = refbit(value, 4) == 1
		ppu.isBgVisible = refbit(value, 3) == 1
	case 0x2003:
		ppu.oamAddress = value
	case 0x2004:
		ppu.oamWriteCount++
		ppu.tmpOam[ppu.oamWriteCount-1] = value
		if ppu.oamWriteCount == 4{
			ppu.oamWriteCount = 0
		}
	case 0x2005:
		switch ppu.writeCount{
		case 0:
			ppu.scrollOffsetX = value
			ppu.writeCount++
		case 1:
			ppu.scrollOffsetY = value
			ppu.writeCount = 0
		}
	case 0x2006:
		switch ppu.writeCount{
		case 0:
			ppu.ppuAddress = uint16(value) * 0x100
			ppu.writeCount++
		case 1:
			ppu.ppuAddress += uint16(value)
			ppu.writeCount = 0
		}
	case 0x2007:
		ppu.ppuAddress += ppu.ppuAddressInc
	case 0x4014:
	}
}

func (ppu Ppu) readPpuRegister(address uint16) uint8{
	switch address{
	case 0x2002:
		var vBlankFlag uint8 = 0
		if ppu.vBlank{
			vBlankFlag = 1
		}
		var spriteHitFlag uint8 = 0
		if ppu.spriteHit{
			spriteHitFlag = 1
		}
		ppu.vBlank = false
		ppu.oamWriteCount = 0
		ppu.writeCount = 0
		return vBlankFlag * 0x80 + spriteHitFlag * 0x40
	case 0x2007:
		return ppu.readVram(address)
	}
	return 0x00
}