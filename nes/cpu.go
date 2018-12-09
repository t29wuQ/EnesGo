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

func (cpu Cpu) readWram(address uint16) uint8{
	return cpu.wram[address]
}

func (cpu Cpu) immediate() uint16{
	address := cpu.programCounter + 1
	cpu.programCounter += 2
	return address
}

func (cpu Cpu) zeropage() uint16{
	address := cpu.programCounter + 1
	cpu.programCounter += 2
	return  uint16(cpu.readWram(address))
}

func (cpu Cpu) zeropageX() uint16{
	 address := cpu.programCounter + 1;
	 cpu.programCounter += 2
	 var eaddress uint8 = uint8(cpu.readWram(address) + cpu.registerX)
	 return uint16(eaddress)
}

func (cpu Cpu) zeropageY() uint16{
	address := cpu.programCounter + 1;
	cpu.programCounter += 2
	var eaddress uint8 = uint8(cpu.readWram(address) + cpu.registerY)
	return uint16(eaddress)
}

func (cpu Cpu) absolute() uint16{
	var address uint16 = uint16(cpu.readWram(cpu.programCounter + 2)) * 0x100 + uint16(cpu.readWram(cpu.programCounter + 1))
	cpu.programCounter += 3
	return address
}

func (cpu Cpu) absoluteX() uint16{
	var address uint16 = uint16(cpu.readWram(cpu.programCounter + 2)) * 0x100 + uint16(cpu.readWram(cpu.programCounter + 1))
	cpu.programCounter += 3
	return address + uint16(cpu.registerX)
}

func (cpu Cpu) absoluteY() uint16{
	var address uint16 = uint16(cpu.readWram(cpu.programCounter + 2)) * 0x100 + uint16(cpu.readWram(cpu.programCounter + 1))
	cpu.programCounter += 3
	return address + uint16(cpu.registerY)
}

func (cpu Cpu) indirectX() uint16{
	var tmp uint8 = uint8(cpu.readWram(cpu.programCounter + 1) + cpu.registerX)
	cpu.programCounter++;
	var address uint16 = uint16(cpu.readWram(uint16(tmp)))
	tmp++
	address |= uint16(cpu.readWram(uint16(tmp))) << 8
	return address
}

func (cpu Cpu) indirectY() uint16{
	tmp := cpu.readWram(cpu.programCounter + 1)
	cpu.programCounter++;
	var address uint16 = uint16(cpu.readWram(uint16(tmp)))
	tmp++
	address |= uint16(cpu.readWram(uint16(tmp))) << 8
	return address + uint16(cpu.registerY)
}