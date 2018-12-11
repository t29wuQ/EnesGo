package nes

func btoi(b bool) int{
	if b {
		return 1
	} else{
		return 0
	}
}

func refbit(i uint8, b uint) int {
    return int((i >> b)) & 1
}


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

func (cpu Cpu) writeWram(address uint16, value uint8){
	cpu.wram[address] = value
}

func (cpu Cpu) readWram(address uint16) uint8{
	return cpu.wram[address]
}

//アドレッシングモードここから

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


//アドレッシングモードここまで


//命令セットここから

func (cpu Cpu) loadRegister(address uint16, register *uint8){
	*register = cpu.readWram(address)
	cpu.setFlagNZ(int(*register))
}

func (cpu Cpu) txs(){
	cpu.programCounter++
	cpu.registerS = cpu.registerX
}

func (cpu Cpu) copyRegister(source *uint8, target *uint8){
	cpu.programCounter++
	*target = *source
	cpu.setFlagNZ(int(*target))
}

func (cpu Cpu) adc(address uint16){
	var sum int = int(cpu.registerA) + int(cpu.readWram(address)) + int(btoi(cpu.cFlag))
	cpu.registerA = uint8(sum)
	cpu.setFlagNZ(int(cpu.registerA))
	cpu.vFlag = ((((cpu.registerA ^ cpu.readWram(address)) & 0x80) != 0) && (((int(cpu.registerA) ^ sum) & 0x80)) != 0);
	cpu.cFlag = sum > 0xff
}

func (cpu Cpu) and(address uint16){
	cpu.registerA &= cpu.readWram(address)
	cpu.setFlagNZ(int(cpu.registerA))
}

func (cpu Cpu) asl(value *uint8){
	cpu.cFlag = (*value >> 7) == 1
	*value <<= 1
	cpu.setFlagNZ(int(*value))
}

func (cpu Cpu) bit(address uint16){
	value := cpu.readWram(address)
	cpu.nFlag = refbit(value, 7) != 0
	cpu.vFlag = refbit(value, 6) != 0
	cpu.zFlag = (cpu.registerA & value) == 0
}

func (cpu Cpu) comparison(address uint16, register *uint8){
	var tmp int = int(*register - cpu.readWram(address))
	cpu.setFlagNZ(int(tmp))
	cpu.cFlag = tmp >= 0
}

func (cpu Cpu) dec(address uint16){
	value := cpu.readWram(address)
	value--
	cpu.writeWram(address, value)
	cpu.setFlagNZ(int(value))
}

func (cpu Cpu) inc(address uint16){
	value := cpu.readWram(address)
	value++
	cpu.writeWram(address, value)
	cpu.setFlagNZ(int(value))
}

func (cpu Cpu) eor(address uint16){
	cpu.registerA ^= cpu.readWram(address)
	cpu.setFlagNZ(int(cpu.registerA))
}

func (cpu Cpu) lsr(value *uint8){
	cpu.cFlag = refbit(*value, 0) != 0
	*value >>= 1
	cpu.setFlagNZ(int(*value))
}

func (cpu Cpu) ora(address uint16){
	cpu.registerA |= cpu.readWram(address)
	cpu.setFlagNZ(int(cpu.registerA))
}

func (cpu Cpu) rol(value *uint8){
	var tmp uint8 = uint8(refbit(*value, 7))
	*value = (*value << 1) + uint8(btoi(cpu.cFlag))
	cpu.cFlag = tmp != 0
	cpu.setFlagNZ(int(*value))
}

func (cpu Cpu) ror(value *uint8){
	var tmp uint8 = uint8(refbit(*value, 0))
	*value = (*value >> 1) + uint8(btoi(cpu.cFlag) * 0x80)
	cpu.cFlag = tmp != 0
	cpu.setFlagNZ(int(*value))
}

func (cpu Cpu) SBC(address uint16){
	cpu.cFlag = cpu.registerA >= cpu.readWram(address) + uint8(btoi(cpu.cFlag) ^ 1)
	var sub uint8 = cpu.registerA - cpu.readWram(address) - uint8(btoi(cpu.cFlag) ^ 1)
	if cpu.registerA < cpu.readWram(address) + uint8(btoi(cpu.cFlag) ^ 1){
		sub++
	}
	cpu.registerA = sub
	cpu.setFlagNZ(int(cpu.registerA))
	cpu.vFlag = (!(((cpu.registerA ^ cpu.readWram(address)) & 0x80) != 0) && (((cpu.registerA ^ sub) & 0x80)) != 0);
}

func (cpu Cpu) push(value uint8){
	cpu.writeWram(0x100 + uint16(cpu.registerS), value)
	cpu.registerS--;
}

func (cpu Cpu) JMP(address uint16){
	cpu.programCounter = address
}

func (cpu Cpu) branch(flag bool){
	if (flag){
		var address int8 = 0
		var tmp uint8 = cpu.readWram(cpu.programCounter + 1)
		if ((tmp >> 7) == 1){
			tmp ^= 0xff
			address = int8(-(tmp + 1))
		}else{
			address = int8(tmp)
		}
		cpu.programCounter = uint16(address) + cpu.programCounter + 2
	} else{
		cpu.programCounter += 2
	}
}



//命令セットここまで

func (cpu Cpu) setFlagNZ(value int){
	
}