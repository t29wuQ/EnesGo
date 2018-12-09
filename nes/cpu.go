package nes



type Cpu struct{
	wram [0x10000]uint8
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

func CreateCpu(prom []uint8) *Cpu{
	cpu := new(Cpu)
	loopCount := int(cpu.programCounter)
	for _, b := range prom {
		cpu.wram[loopCount] = b
		loopCount++;
	}
	cpu.programCounter = 0x8000
	cpu.registerS = 0xff
	return cpu
}