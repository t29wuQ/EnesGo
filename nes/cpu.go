package nes

func btoi(b bool) int{
	if b {
		return 1
	} else{
		return 0
	}
}

func refbit(i uint8, b uint) uint8 {
    return (i >> b) & 1
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

var cpuCycle [0x100]int = [0x100]int{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 4, 2, 5, 5, 4, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 3, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7}

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

func (cpu Cpu) getRegisterP() uint8{
	return uint8(btoi(cpu.nFlag) * 0x80 +
	btoi(cpu.vFlag) * 0x40 + 0x20 +
	btoi(cpu.bFlag) * 0x10 +
	btoi(cpu.iFlag) * 0x04 + 
	btoi(cpu.zFlag) * 0x02 + 
	btoi(cpu.cFlag))
}

func (cpu Cpu) setRegisterP(value uint8){
	cpu.nFlag = refbit(value, 7) != 0
	cpu.vFlag = refbit(value, 6) != 0
	cpu.bFlag = refbit(value, 4) != 0
	cpu.iFlag = refbit(value, 2) != 0
	cpu.zFlag = refbit(value, 1) != 0
	cpu.cFlag = refbit(value, 0) != 0
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
	cpu.setFlagNZ(*register)
}

func (cpu Cpu) txs(){
	cpu.programCounter++
	cpu.registerS = cpu.registerX
}

func (cpu Cpu) copyRegister(source *uint8, target *uint8){
	cpu.programCounter++
	*target = *source
	cpu.setFlagNZ(*target)
}

func (cpu Cpu) adc(address uint16){
	var sum int = int(cpu.registerA) + int(cpu.readWram(address)) + int(btoi(cpu.cFlag))
	cpu.registerA = uint8(sum)
	cpu.setFlagNZ(cpu.registerA)
	cpu.vFlag = ((((cpu.registerA ^ cpu.readWram(address)) & 0x80) != 0) && (((int(cpu.registerA) ^ sum) & 0x80)) != 0);
	cpu.cFlag = sum > 0xff
}

func (cpu Cpu) and(address uint16){
	cpu.registerA &= cpu.readWram(address)
	cpu.setFlagNZ(cpu.registerA)
}

func (cpu Cpu) asl(value *uint8){
	cpu.cFlag = (*value >> 7) == 1
	*value <<= 1
	cpu.setFlagNZ(*value)
}

func (cpu Cpu) bit(address uint16){
	value := cpu.readWram(address)
	cpu.nFlag = refbit(value, 7) != 0
	cpu.vFlag = refbit(value, 6) != 0
	cpu.zFlag = (cpu.registerA & value) == 0
}

func (cpu Cpu) comparison(address uint16, register *uint8){
	var tmp int = int(*register - cpu.readWram(address))
	cpu.setFlagNZ(uint8(tmp))
	cpu.cFlag = tmp >= 0
}

func (cpu Cpu) dec(address uint16){
	value := cpu.readWram(address)
	value--
	cpu.writeWram(address, value)
	cpu.setFlagNZ(value)
}

func (cpu Cpu) inc(address uint16){
	value := cpu.readWram(address)
	value++
	cpu.writeWram(address, value)
	cpu.setFlagNZ(value)
}

func (cpu Cpu) eor(address uint16){
	cpu.registerA ^= cpu.readWram(address)
	cpu.setFlagNZ(cpu.registerA)
}

func (cpu Cpu) lsr(value *uint8){
	cpu.cFlag = refbit(*value, 0) != 0
	*value >>= 1
	cpu.setFlagNZ(*value)
}

func (cpu Cpu) ora(address uint16){
	cpu.registerA |= cpu.readWram(address)
	cpu.setFlagNZ(cpu.registerA)
}

func (cpu Cpu) rol(value *uint8){
	var tmp uint8 = uint8(refbit(*value, 7))
	*value = (*value << 1) + uint8(btoi(cpu.cFlag))
	cpu.cFlag = tmp != 0
	cpu.setFlagNZ(*value)
}

func (cpu Cpu) ror(value *uint8){
	var tmp uint8 = uint8(refbit(*value, 0))
	*value = (*value >> 1) + uint8(btoi(cpu.cFlag) * 0x80)
	cpu.cFlag = tmp != 0
	cpu.setFlagNZ(*value)
}

func (cpu Cpu) sbc(address uint16){
	cpu.cFlag = cpu.registerA >= cpu.readWram(address) + uint8(btoi(cpu.cFlag) ^ 1)
	var sub uint8 = cpu.registerA - cpu.readWram(address) - uint8(btoi(cpu.cFlag) ^ 1)
	if cpu.registerA < cpu.readWram(address) + uint8(btoi(cpu.cFlag) ^ 1){
		sub++
	}
	cpu.registerA = sub
	cpu.setFlagNZ(cpu.registerA)
	cpu.vFlag = (!(((cpu.registerA ^ cpu.readWram(address)) & 0x80) != 0) && (((cpu.registerA ^ sub) & 0x80)) != 0);
}

func (cpu Cpu) push(value uint8){
	cpu.writeWram(0x100 + uint16(cpu.registerS), value)
	cpu.registerS--;
}

func (cpu Cpu) jmp(address uint16){
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

func (cpu Cpu) setFlagNZ(value uint8){
	cpu.zFlag = value == 0
	cpu.nFlag = refbit(value, 7) != 0
}

func (cpu Cpu) clearFlag(flag *bool){
	*flag = false
	cpu.programCounter++
}

func (cpu Cpu) setFlag(flag *bool){
	*flag = true
	cpu.programCounter++
}

//命令セットここまで

func (cpu Cpu) Execute() int{
	opcode := cpu.readWram(cpu.programCounter)
	switch opcode{
	case 0xA9:
		cpu.loadRegister(cpu.immediate(), &cpu.registerA);
	case 0xA5:
		cpu.loadRegister(cpu.zeropage(), &cpu.registerA);
	case 0xB5:
		cpu.loadRegister(cpu.zeropageX(), &cpu.registerA);
	case 0xAD:
		cpu.loadRegister(cpu.absolute(), &cpu.registerA);
	case 0xBD:
		cpu.loadRegister(cpu.absoluteX(), &cpu.registerA);
	case 0xB9:
		cpu.loadRegister(cpu.absoluteY(), &cpu.registerA);
	case 0xA1:
		cpu.loadRegister(cpu.indirectX(), &cpu.registerA);
	case 0xB1:
		cpu.loadRegister(cpu.indirectY(), &cpu.registerA);
	case 0xA2:
		cpu.loadRegister(cpu.immediate(), &cpu.registerX);
	case 0xA6:
		cpu.loadRegister(cpu.zeropage(), &cpu.registerX);
	case 0xB6:
		cpu.loadRegister(cpu.zeropageY(), &cpu.registerX);
	case 0xAE:
		cpu.loadRegister(cpu.absolute(), &cpu.registerX);
	case 0xBE:
		cpu.loadRegister(cpu.absoluteY(), &cpu.registerX);
	case 0xA0:
		cpu.loadRegister(cpu.immediate(), &cpu.registerY);
	case 0xA4:
		cpu.loadRegister(cpu.zeropage(), &cpu.registerY);
	case 0xB4:
		cpu.loadRegister(cpu.zeropageX(), &cpu.registerY);
	case 0xAC:
		cpu.loadRegister(cpu.absolute(), &cpu.registerY);
	case 0xBC:
		cpu.loadRegister(cpu.absoluteX(), &cpu.registerY);
	case 0x85:
		cpu.writeWram(cpu.zeropage(), cpu.registerA);
	case 0x95:
		cpu.writeWram(cpu.zeropageX(), cpu.registerA);
	case 0x8D:
		cpu.writeWram(cpu.absolute(), cpu.registerA);
	case 0x9D:
		cpu.writeWram(cpu.absoluteX(), cpu.registerA);
	case 0x99:
		cpu.writeWram(cpu.absoluteY(), cpu.registerA);
	case 0x81:
		cpu.writeWram(cpu.indirectX(), cpu.registerA);
	case 0x91:
		cpu.writeWram(cpu.indirectY(), cpu.registerA);
	case 0x86:
		cpu.writeWram(cpu.zeropage(), cpu.registerX);
	case 0x96:
		cpu.writeWram(cpu.zeropageY(), cpu.registerX);
	case 0x8E:
		cpu.writeWram(cpu.absolute(), cpu.registerX);
	case 0x84:
		cpu.writeWram(cpu.zeropage(), cpu.registerY);
	case 0x94:
		cpu.writeWram(cpu.zeropageX(), cpu.registerY);
	case 0x8C:
		cpu.writeWram(cpu.absolute(), cpu.registerY);
	case 0xAA:
		cpu.copyRegister(&cpu.registerA, &cpu.registerX);
	case 0xA8:
		cpu.copyRegister(&cpu.registerA, &cpu.registerY);
	case 0xBA:
		cpu.copyRegister(&cpu.registerS, &cpu.registerX);
	case 0x8A:
		cpu.copyRegister(&cpu.registerX, &cpu.registerA);
	case 0x9A:
		cpu.txs();
	case 0x98:
		cpu.copyRegister(&cpu.registerY, &cpu.registerA);
	case 0x69:
		cpu.adc(cpu.immediate());
	case 0x65:
		cpu.adc(cpu.zeropage());
	case 0x75:
		cpu.adc(cpu.zeropageX());
	case 0x6D:
		cpu.adc(cpu.absolute());
	case 0x7D:
		cpu.adc(cpu.absoluteX());
	case 0x79:
		cpu.adc(cpu.absoluteY());
	case 0x61:
		cpu.adc(cpu.indirectX());
	case 0x71:
		cpu.adc(cpu.indirectY());
	case 0x29:
		cpu.and(cpu.immediate());
	case 0x25:
		cpu.and(cpu.zeropage());
	case 0x35:
		cpu.and(cpu.zeropageX());
	case 0x2D:
		cpu.and(cpu.absolute());
	case 0x3D:
		cpu.and(cpu.absoluteX());
	case 0x39:
		cpu.and(cpu.absoluteY());
	case 0x21:
		cpu.and(cpu.indirectX());
	case 0x31:
		cpu.and(cpu.indirectY());
	case 0x0A:
		cpu.asl(&cpu.registerA);
		cpu.programCounter++;
	case 0x06:
		cpu.asl(&cpu.wram[cpu.zeropage()]);
	case 0x16:
		cpu.asl(&cpu.wram[cpu.zeropageX()]);
	case 0x0E:
		cpu.asl(&cpu.wram[cpu.absolute()]);
	case 0x1E:
		cpu.asl(&cpu.wram[cpu.absoluteX()]);
	case 0x24:
		cpu.bit(cpu.zeropage());
	case 0x2C:
		cpu.bit(cpu.absolute());
	case 0xC9:
		cpu.comparison(cpu.immediate(), &cpu.registerA);
	case 0xC5:
		cpu.comparison(cpu.zeropage(), &cpu.registerA);
	case 0xD5:
		cpu.comparison(cpu.zeropageX(), &cpu.registerA);
	case 0xCD:
		cpu.comparison(cpu.absolute(), &cpu.registerA);
	case 0xDD:
		cpu.comparison(cpu.absoluteX(), &cpu.registerA);
	case 0xD9:
		cpu.comparison(cpu.absoluteY(), &cpu.registerA);
	case 0xC1:
		cpu.comparison(cpu.indirectX(), &cpu.registerA);
	case 0xD1:
		cpu.comparison(cpu.indirectY(), &cpu.registerA);
	case 0xE0:
		cpu.comparison(cpu.immediate(), &cpu.registerX);
	case 0xE4:
		cpu.comparison(cpu.zeropage(), &cpu.registerX);
	case 0xEC:
		cpu.comparison(cpu.absolute(), &cpu.registerX);
	case 0xC0:
		cpu.comparison(cpu.immediate(), &cpu.registerY);
	case 0xC4:
		cpu.comparison(cpu.zeropage(), &cpu.registerY);
	case 0xCC:
		cpu.comparison(cpu.absolute(), &cpu.registerY);
	case 0xC6:
		cpu.dec(cpu.zeropage());
	case 0xD6:
		cpu.dec(cpu.zeropageX());
	case 0xCE:
		cpu.dec(cpu.absolute());
	case 0xDE:
		cpu.dec(cpu.absoluteX());
	/*
	 * Xをデクリメント
	 * N: 演算結果の最上位ビット
	 * Z: 演算結果が0であるか
	 */
	case 0xCA:
		cpu.programCounter++;
		cpu.registerX--;
		cpu.setFlagNZ(cpu.registerX);
	/*
	 * Yをデクリメント
	 * N: 演算結果の最上位ビット
	 * Z: 演算結果が0であるか
	 */
	case 0x88:
		cpu.programCounter++;
		cpu.registerY--;
		cpu.setFlagNZ(cpu.registerY);
	case 0x49:
		cpu.eor(cpu.immediate());
	case 0x45:
		cpu.eor(cpu.zeropage());
	case 0x55:
		cpu.eor(cpu.zeropageX());
	case 0x4D:
		cpu.eor(cpu.absolute());
	case 0x5D:
		cpu.eor(cpu.absoluteX());
	case 0x59:
		cpu.eor(cpu.absoluteY());
	case 0x41:
		cpu.eor(cpu.indirectX());
	case 0x51:
		cpu.eor(cpu.indirectY());
	case 0xE6:
		cpu.inc(cpu.zeropage());
	case 0xF6:
		cpu.inc(cpu.zeropageX());
	case 0xEE:
		cpu.inc(cpu.absolute());
	case 0xFE:
		cpu.inc(cpu.absoluteX());
	/*
	 * Xをインクリメント
	 * N: 演算結果の最上位ビット
	 * Z: 演算結果が0であるか
	 */
	case 0xE8:
		cpu.programCounter++;
		cpu.registerX++;
		cpu.setFlagNZ(cpu.registerX);
	/*
	 * Yをインクリメント
	 * N: 演算結果の最上位ビット
	 * Z: 演算結果が0であるか
	 */
	case 0xC8:
		cpu.programCounter++;
		cpu.registerY++;
		cpu.setFlagNZ(cpu.registerY);
	case 0x4A:
		cpu.lsr(&cpu.registerA);
		cpu.programCounter++;
	case 0x46:
		cpu.lsr(&cpu.wram[cpu.zeropage()]);
	case 0x56:
		cpu.lsr(&cpu.wram[cpu.zeropageX()]);
	case 0x4E:
		cpu.lsr(&cpu.wram[cpu.absolute()]);
	case 0x5E:
		cpu.lsr(&cpu.wram[cpu.absoluteX()]);
	case 0x09:
		cpu.ora(cpu.immediate());
	case 0x05:
		cpu.ora(cpu.zeropage());
	case 0x15:
		cpu.ora(cpu.zeropageX());
	case 0x0D:
		cpu.ora(cpu.absolute());
	case 0x1D:
		cpu.ora(cpu.absoluteX());
	case 0x19:
		cpu.ora(cpu.absoluteY());
	case 0x01:
		cpu.ora(cpu.indirectX());
	case 0x11:
		cpu.ora(cpu.indirectY());
	case 0x2A:
		cpu.rol(&cpu.registerA);
		cpu.programCounter++;
	case 0x26:
		cpu.rol(&cpu.wram[cpu.zeropage()]);
	case 0x36:
		cpu.rol(&cpu.wram[cpu.zeropageX()]);
	case 0x2E:
		cpu.rol(&cpu.wram[cpu.absolute()]);
	case 0x3E:
		cpu.rol(&cpu.wram[cpu.absoluteX()]);
	case 0x6A:
		cpu.ror(&cpu.registerA);
		cpu.programCounter++;
	case 0x66:
		cpu.ror(&cpu.wram[cpu.zeropage()]);
	case 0x76:
		cpu.ror(&cpu.wram[cpu.zeropageX()]);
	case 0x6E:
		cpu.ror(&cpu.wram[cpu.absolute()]);
	case 0x7E:
		cpu.ror(&cpu.wram[cpu.absoluteX()]);
	case 0xE9:
		cpu.sbc(cpu.immediate());
	case 0xE5:
		cpu.sbc(cpu.zeropage());
	case 0xF5:
		cpu.sbc(cpu.zeropageX());
	case 0xED:
		cpu.sbc(cpu.absolute());
	case 0xFD:
		cpu.sbc(cpu.absoluteX());
	case 0xF9:
		cpu.sbc(cpu.absoluteY());
	case 0xE1:
		cpu.sbc(cpu.indirectX());
	case 0xF1:
		cpu.sbc(cpu.indirectY());
	case 0x48:
		cpu.push(cpu.registerA);
		cpu.programCounter++;
	case 0x08:
		cpu.push(cpu.getRegisterP());
		cpu.programCounter++;
	/*
	 * PLA
	 * スタックからAにポップアップ
	 * N: POPした値の最上位ビット
	 * Z: POPした値が0であるか
	 */
	case 0x68:
		cpu.programCounter++;
		cpu.registerS++;
		cpu.registerA = cpu.readWram(0x0100 + uint16(cpu.registerS));
		cpu.setFlagNZ(cpu.registerA);
	/*
	 * スタックからPにポップアップ
	 */
	case 0x28:
		cpu.programCounter++;
		cpu.registerS++;
		cpu.setRegisterP(cpu.readWram(0x0100 + uint16(cpu.registerS)))
	case 0x4C:
		cpu.jmp(cpu.absolute());
	case 0x6C:
		cpu.jmp(uint16(cpu.readWram(uint16(cpu.readWram(cpu.programCounter + 1)) + 1)) * 0x100 + uint16(cpu.readWram(uint16(cpu.readWram(cpu.programCounter + 1)))))
	/*
	 * サブルーチンを呼び出す
	 * 元のPCを上位, 下位バイトの順にcpu.pushする
	 * この時保存するPCはJSRの最後のバイトアドレス
	 * JSR
	 */
	case 0x20:
		savePC := (uint16)(cpu.programCounter + 2);
		cpu.push((byte)(savePC >> 8));
		cpu.push((byte)((savePC << 8) >> 8));
		cpu.programCounter = cpu.absolute();
	/*
	 * サブルーチンから復帰する
	 * 下位, 上位バイトの順にpopする
	 * RTS
	 */
	case 0x60:
		cpu.programCounter = uint16(cpu.readWram(0x0100 + uint16(cpu.registerS) + 2)) * 0x0100 + uint16(cpu.readWram(0x0100 + uint16(cpu.registerS) + 1))
		cpu.registerS += 2;
		cpu.programCounter++;
	/*
	 * 割り込みハンドラから復帰
	 * RTI
	 */
	case 0x40:
		cpu.registerS++;
		cpu.setRegisterP(cpu.readWram(0x0100 + uint16(cpu.registerS)));
		cpu.programCounter = uint16(cpu.readWram(0x0100 + uint16(cpu.registerS) + 2)) * 0x0100 + uint16(cpu.readWram(0x0100 + uint16(cpu.registerS) + 1))
		cpu.registerS += 2;
	case 0x90:
		cpu.branch(!cpu.cFlag);
	case 0xB0:
		cpu.branch(cpu.cFlag);
	case 0xF0:
		cpu.branch(cpu.zFlag);
	case 0x30:
		cpu.branch(cpu.nFlag);
	case 0xD0:
		cpu.branch(!cpu.zFlag);
	case 0x10:
		cpu.branch(!cpu.nFlag);
	case 0x50:
		cpu.branch(!cpu.vFlag);
	case 0x70:
		cpu.branch(cpu.vFlag);
	case 0x18:
		cpu.clearFlag(&cpu.cFlag);
	case 0xD8:
		cpu.programCounter++;
	case 0x58:
		cpu.clearFlag(&cpu.iFlag);
	case 0xB8:
		cpu.clearFlag(&cpu.vFlag);
	case 0x38:
		cpu.setFlag(&cpu.cFlag);
	case 0x78:
		cpu.setFlag(&cpu.iFlag);
	/*
	 * ソフトウェア割り込みを起こす
	 * BRK
	 */
	case 0x00:
		cpu.push((byte)(cpu.programCounter >> 8));
		cpu.push((byte)((cpu.programCounter << 8) >> 8));
		cpu.push(cpu.getRegisterP());
		cpu.programCounter = uint16(cpu.readWram(0xFFFF)) * 0x100 + uint16(cpu.readWram(0xFFFE));
	/*
	 * 空の命令を実効
	 * NOP
	 */
	case 0xEA:
		cpu.programCounter++;
	}
	return cpuCycle[opcode]
}