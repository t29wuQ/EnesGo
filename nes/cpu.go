package nes



type Cpu struct{
	Wram []uint8
	registerA uint8
	registerX uint8
	registerY uint8
	registerS uint8
	programCounter uint16
	nFlag bool
	vFlag bool
	bFlag bool
	iFlag bool
	zFlag bool
	cFlag bool
}