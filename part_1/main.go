package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

var instTable = []inst{
	inst{
		opType: "mov", 
		opCode: 0b100010,
		opCodeDiff: 2,
		params: []P{
			P{c: 0b00000010, n: "d", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
	inst{
		opType: "movImm", 
		opCode: 0b1011,
		opCodeDiff: 4,
		params: []P{
			P{c: 0b00001000, n: "w", len: 3},
			P{c: 0b00000111, n: "reg", len: 0},
		},
	},
	inst{
		opType: "add", 
		opCode: 0b000000,
		opCodeDiff: 2,
		params: []P{
			P{c: 0b00000010, n: "d", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
	inst{
		opType: "addImmReg", 
		opCode: 0b1000,
		opCodeDiff: 4,
		params: []P{
			P{c: 0b00000010, n: "s", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
	inst{
		opType: "addImmAcc", 
		opCode: 0b0000010,
		opCodeDiff: 1,
		params: []P{
			P{c: 0b00000001, n: "w", len: 0},
		},
	},
	inst{
		opType: "sub", 
		opCode: 0b001010,
		opCodeDiff: 2,
		params: []P{
			P{c: 0b00000010, n: "d", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
	inst{
		opType: "subImmReg", 
		opCode: 0b1000,
		opCodeDiff: 4,
		params: []P{
			P{c: 0b00000010, n: "s", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
	inst{
		opType: "subImmAcc", 
		opCode: 0b0010110,
		opCodeDiff: 1,
		params: []P{
			P{c: 0b00000001, n: "w", len: 0},
		},
	},
	inst{
		opType: "cmp", 
		opCode: 0b001110,
		opCodeDiff: 2,
		params: []P{
			P{c: 0b00000010, n: "d", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
	inst{
		opType: "cmpImmReg", 
		opCode: 0b100000,
		opCodeDiff: 4,
		params: []P{
			P{c: 0b00000010, n: "s", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
	inst{
		opType: "cmpImmAcc", 
		opCode: 0b0011110,
		opCodeDiff: 1,
		params: []P{
			P{c: 0b00000001, n: "w", len: 0},
		},
	},
} 

var ax []uint8 = []uint8{0b00000000, 0b00000000}
var bx []uint8 = []uint8{0b00000000, 0b00000000}
var cx []uint8 = []uint8{0b00000000, 0b00000000}
var dx []uint8 = []uint8{0b00000000, 0b00000000}
var sp uint16 = 0b0000000000000000
var bp uint16 = 0b0000000000000000
var si uint16 = 0b0000000000000000
var di uint16 = 0b0000000000000000

func (inst inst) getP(value string) P {
	for index := 0; index < len(inst.params); index++ {
		if inst.params[index].n == value {
			return inst.params[index]
		}
	}
	panic("no such param for inst")
}

func main() {
	//data, err := ioutil.ReadFile("add_sub_cmp")
	data, err := ioutil.ReadFile("add_flags")
	fmt.Printf("inst %08b\n", data)
	check(err)
	var max = 30
	//for len(data) > 0 && max > 0 {
	for len(data) > 0 {
		data = disassembleAndReturn(data)
		max -= 1
		//fmt.Printf("inst %08b\n", data)
	}
	printRegisters()
}

func printRegisters() {
	var result uint16 = uint16(ax[1]) << 8 | uint16(ax[0]) 
	fmt.Printf("%s = %08b, %08b => %d\n", "ax", ax[0], ax[1], result)
	result = uint16(bx[1]) << 8 | uint16(bx[0]) 
	fmt.Printf("%s = %08b, %08b => %d\n", "bx", bx[0], bx[1], result)
	result = uint16(cx[1]) << 8 | uint16(cx[0]) 
	fmt.Printf("%s = %08b, %08b => %d\n", "cx", cx[0], cx[1], result)
	result = uint16(dx[1]) << 8 | uint16(dx[0]) 
	fmt.Printf("%s = %08b, %08b => %d\n", "dx", dx[0], dx[1], result)
	fmt.Printf("%s = %016b => %d\n", "sp", sp, sp)
	fmt.Printf("%s = %016b => %d\n", "bp", bp, bp)
	fmt.Printf("%s = %016b => %d\n", "si", si, si)
	fmt.Printf("%s = %016b => %d\n", "di", di, di)
}

func disassembleAndReturn(instructions []uint8) []uint8 {
	for i := 0; i < len(instTable); i++ {
		var inst = instTable[i]
		var opCode = instructions[0] >> inst.opCodeDiff
		//fmt.Printf("opCode %08b\n", instructions[0])
		if (opCode == inst.opCode) {
			if inst.opType == "mov" {
				return disMovComplex(instructions, inst)
			} 
			if inst.opType == "movImm" {
				return disMovImmediateToRegister(instructions, inst)
			}
			if inst.opType == "add" {
				return disAddComplex("add", instructions, inst)
			}
			if inst.opType == "addImmReg" {
				return disAddImmediateToRegister("add", instructions, inst)
			}
			if inst.opType == "addImmAcc" {
				return disAddAccumulatorToRegister("add", instructions, inst)
			}
			if inst.opType == "sub" {
				return disAddComplex("sub", instructions, inst)
			}
			if inst.opType == "subImmReg" {
				return disAddImmediateToRegister("sub", instructions, inst)
			}
			if inst.opType == "subImmAcc" {
				return disAddAccumulatorToRegister("sub", instructions, inst)
			}
			if inst.opType == "cmp" {
				return disAddComplex("cmp", instructions, inst)
			}
			if inst.opType == "cmpImmReg" {
				return disAddImmediateToRegister("cmp", instructions, inst)
			}
			if inst.opType == "cmpImmAcc" {
				return disAddAccumulatorToRegister("cmp", instructions, inst)
			}
		}
		//else {
		//	return disassembleJumps(instructions)
		//}
	}
	return disassembleJumps(instructions)
	//return instructions
}

func disassembleJumps(instructions []uint8) []uint8 {
	var removeCount = 0
	var instruction = instructions[0]
	removeCount += 1
	var name = ""
	switch (instruction) {
		case 0b01110100:
			name = "je"
		case 0b01111100:
			name = "jl"
		case 0b01111110:
			name = "jle"
		case 0b01110010:
			name = "jb"
		case 0b01110110:
			name = "jbe"
		case 0b01111010:
			name = "jp"
		case 0b01110000:
			name = "jo"
		case 0b01111000:
			name = "js"
		case 0b01110101:
			name = "jne"
		case 0b01111101:
			name = "jnl"
		case 0b01111111:
			name = "jg"
		case 0b01110011:
			name = "jnb"
		case 0b01110111:
			name = "ja"
		case 0b01111011:
			name = "jnp"
		case 0b01110001:
			name = "jno"
		case 0b01111001:
			name = "jns"
		case 0b11100010:
			name = "loop"
		case 0b11100001:
			name = "loopz"
		case 0b11100000:
			name = "loopnz"
		case 0b11100011:
			name = "jcxz"
	}
	if (name == "") {
		return instructions
	}
	removeCount += 1
	fmt.Printf("%s %d \n", name, instructions[1])

	return instructions[removeCount:]
}

func updateRegister(reg string, value uint8) {
	switch(reg) {
	case "al":
		ax[0] = value 
	case "ah":
		ax[1] = value 
	}
}

func updateRegisterToRegister(destination string, source string) {
	var src = getRegister(source)
	switch(destination) {
	case "ax":
		ax[0] = src[0] 
		ax[1] = src[1] 
	case "bx":
		bx[0] = src[0] 
		bx[1] = src[1] 
	case "cx":
		cx[0] = src[0] 
		cx[1] = src[1] 
	case "dx":
		dx[0] = src[0] 
		dx[1] = src[1]
	case "sp":
		sp = u8Tou16(src[0], src[1])
	case "bp":
		bp = u8Tou16(src[0], src[1])
	case "si":
		si = u8Tou16(src[0], src[1])
	case "di":
		di = u8Tou16(src[0], src[1])
	}
}

func getRegister(reg string) []uint8 {
	switch(reg) {
	case "ax":
		return ax
	case "bx":
		return bx
	case "cx":
		return cx
	case "dx":
		return dx
	case "sp":
	return u16toArray(sp)
	case "bp":
	return u16toArray(bp)
	case "si":
	return u16toArray(si)
	case "di":
	return u16toArray(di)
	}
	return []uint8{}
}

func u16toArray(input uint16) []uint8 {
	var value1 uint8 = uint8(input)  
	var value2 uint8 = uint8(input >> 8) 
	return []uint8{value1, value2}
}

func u8Tou16(value uint8, value2 uint8) uint16{
	return uint16(value2) << 8 | uint16(value) 
}

func updateRegisterX(reg string, value uint8, value2 uint8) {
	switch(reg) {
	case "ax":
		ax[0] = value 
		ax[1] = value2
	case "bx":
		bx[0] = value 
		bx[1] = value2
	case "cx":
		cx[0] = value 
		cx[1] = value2
	case "dx":
		dx[0] = value 
		dx[1] = value2
	}
	var result uint16 = uint16(value2) << 8 | uint16(value) 
	updateRegister16(reg, result)
}

func updateRegister16(reg string, value uint16) {
	switch(reg) {
	case "sp":
		sp = value
	case "bp":
		bp = value
	case "si":
		si = value
	case "di":
		di = value
	}
}


//only one byte big operation
func disMovImmediateToRegister(instructions []uint8, inst inst) []uint8 {
	var removeCount = 0
	// first one is for opCode
	var instruction = instructions[0]
	var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)
	var reg = clearBits(instruction, inst.getP("reg").c, inst.getP("reg").len)

	var short, wide = transformRegToString(reg)

	if (isBitTrue(w)){
		removeCount += 2
		var result uint16 = uint16(instructions[2]) << 8 | uint16(instructions[1]) 
		updateRegisterX(wide, instructions[1], instructions[2])
		fmt.Printf("mov %s %d\n", wide, result)
	} else {
		removeCount += 1
		updateRegister(short, instructions[1])
		fmt.Printf("mov %s %d\n", short, instructions[1])
	}

	removeCount += 1
	return instructions[removeCount:]
}

func disMovComplex(instructions []uint8, inst inst) []uint8 {
	var removeCount = 0
	// first one is for opCode
	var instruction = instructions[0]
	removeCount += 1
	//every move should have a register part?
	var register = instructions[1]
	removeCount += 1

	//get all the regular values
	var d = clearBits(instruction, inst.getP("d").c, inst.getP("d").len)
	var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)
	var mod = clearBits(register, inst.getP("mod").c, inst.getP("mod").len)
	var reg = clearBits(register, inst.getP("reg").c, inst.getP("reg").len)
	var rm = clearBits(register, inst.getP("rm").c, inst.getP("rm").len)

	//hangle reg to reg
	if (mod == 0b11) {
		//register to register
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		var rm1, rm2 = transformRegToString(rm)
		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = rm2
			} else {
				destination = reg1 
				source = rm1
			}
		} else {
			if (isWide) {
				destination = rm2 
				source = reg2
			} else {
				destination = rm1 
				source = reg1
			}
		}
		updateRegisterToRegister(destination, source)
		fmt.Printf("mov %s %s \n", destination, source)
	}

	// handle direct addressing with no displacement
	if (mod == 0b00) {
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		//var rm1, rm2 = transformRegToString(rm)
		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = "[" + translateMemoryDisplacement00(rm) + "]"
			} else {
				destination = reg1 
				source = "[" + translateMemoryDisplacement00(rm) + "]"
			}
		} else {
			if (isWide) {
				destination = "[" + translateMemoryDisplacement00(rm) + "]"
				source = reg2
			} else {
				destination = "[" + translateMemoryDisplacement00(rm) + "]"
				source = reg1
			}
		}
		fmt.Printf("mov %s, %s \n", destination, source)
	}
	if (mod == 0b01) {
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		//var rm1, rm2 = transformRegToString(rm)
		var result = ""
		var translated = translateMemoryDisplacement00(rm)
		if (translated != "bp") {
			result = " + " + strconv.Itoa(int(instructions[2]))
		}

		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) +  "]"
			} else {
				destination = reg1 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
			}
		} else {
			if (isWide) {
				destination = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
				source = reg2
			} else {
				destination = "[" + translateMemoryDisplacement00(rm) + string(result) + "]"
				source = reg1
			}
		}
		fmt.Printf("mov %s, %s \n", destination, source )
		removeCount += 1
	}
	if (mod == 0b10) {
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		//var rm1, rm2 = transformRegToString(rm)
		var uintResult uint16 = uint16(instructions[3]) << 8 | uint16(instructions[2]) 
		var result = ""
		var translated = translateMemoryDisplacement00(rm)
		if (translated != "bp") {
			result = " + " + strconv.Itoa(int(uintResult))
		}
		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
			} else {
				destination = reg1 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
			}
		} else {
			if (isWide) {
				destination = "[" + translateMemoryDisplacement00(rm) + string(result) +  "]"
				source = reg2
			} else {
				destination = "[" + translateMemoryDisplacement00(rm) + string(result) +  "]"
				source = reg1
			}
		}
		fmt.Printf("mov %s %s \n", destination, source)
		removeCount += 2
	}

	return instructions[removeCount:]
}

func disAddComplex(instrcutionType string, instructions []uint8, inst inst) []uint8 {
	var removeCount = 0
	// first one is for opCode
	var instruction = instructions[0]
	removeCount += 1
	//every move should have a register part?
	var register = instructions[1]
	removeCount += 1

	//get all the regular values
	var d = clearBits(instruction, inst.getP("d").c, inst.getP("d").len)
	var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)
	var mod = clearBits(register, inst.getP("mod").c, inst.getP("mod").len)
	var reg = clearBits(register, inst.getP("reg").c, inst.getP("reg").len)
	var rm = clearBits(register, inst.getP("rm").c, inst.getP("rm").len)

	//hangle reg to reg
	if (mod == 0b11) {
		//register to register
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		var rm1, rm2 = transformRegToString(rm)
		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = rm2
			} else {
				destination = reg1 
				source = rm1
			}
		} else {
			if (isWide) {
				destination = rm2 
				source = reg2
			} else {
				destination = rm1 
				source = reg1
			}
		}
		if (instrcutionType == "sub") {
			var destinationValue = getRegister(destination)
			var sourceValue = getRegister(source)
			var destu16 = u8Tou16(destinationValue[0], destinationValue[1])
			var srcu16 = u8Tou16(sourceValue[0], sourceValue[1])
			updateRegister16(destination, uint16(destu16 - srcu16))
		}
		fmt.Printf("%s %s %s \n", instrcutionType, destination, source)
	}

	// handle direct addressing with no displacement
	if (mod == 0b00) {
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		//var rm1, rm2 = transformRegToString(rm)
		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = "[" + translateMemoryDisplacement00(rm) + "]"
			} else {
				destination = reg1 
				source = "[" + translateMemoryDisplacement00(rm) + "]"
			}
		} else {
			if (isWide) {
				destination = "[" + translateMemoryDisplacement00(rm) + "]"
				source = reg2
			} else {
				destination = "[" + translateMemoryDisplacement00(rm) + "]"
				source = reg1
			}
		}
		fmt.Printf("%s %s, %s \n", instrcutionType, destination, source)
	}
	if (mod == 0b01) {
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		//var rm1, rm2 = transformRegToString(rm)
		var result = ""
		var translated = translateMemoryDisplacement00(rm)
		if (translated != "bp") {
			result = " + " + strconv.Itoa(int(instructions[2]))
		}

		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) +  "]"
			} else {
				destination = reg1 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
			}
		} else {
			if (isWide) {
				destination = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
				source = reg2
			} else {
				destination = "[" + translateMemoryDisplacement00(rm) + string(result) + "]"
				source = reg1
			}
		}
		fmt.Printf("%s %s, %s \n", instrcutionType, destination, source )
		removeCount += 1
	}
	if (mod == 0b10) {
		var destination = ""
		var source = ""
		var isDestinationReg = isBitTrue(d)
		var isWide = isBitTrue(w)
		var reg1, reg2 = transformRegToString(reg)
		//var rm1, rm2 = transformRegToString(rm)
		var uintResult uint16 = uint16(instructions[3]) << 8 | uint16(instructions[2]) 
		var result = ""
		var translated = translateMemoryDisplacement00(rm)
		if (translated != "bp") {
			result = " + " + strconv.Itoa(int(uintResult))
		}
		if (isDestinationReg) {
			if (isWide) {
				destination = reg2 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
			} else {
				destination = reg1 
				source = "[" + translateMemoryDisplacement00(rm) +  string(result) + "]"
			}
		} else {
			if (isWide) {
				destination = "[" + translateMemoryDisplacement00(rm) + string(result) +  "]"
				source = reg2
			} else {
				destination = "[" + translateMemoryDisplacement00(rm) + string(result) +  "]"
				source = reg1
			}
		}
		fmt.Printf("%s %s %s \n", instrcutionType, destination, source)
		removeCount += 2
	}

	return instructions[removeCount:]
}

func disAddImmediateToRegister(instructionType string, instructions []uint8, inst inst) []uint8 {
	var removeCount = 0
	// first one is for opCode
	//var instruction = instructions[0]
	//var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)
	//var s = clearBits(instruction, inst.getP("s").c, inst.getP("s").len)
	removeCount +=1

	var registers = instructions[1]
	removeCount += 1
	var mod = clearBits(registers, inst.getP("mod").c, inst.getP("mod").len)
	var reg = clearBits(registers, inst.getP("reg").c, inst.getP("reg").len)
	var rm = clearBits(registers, inst.getP("rm").c, inst.getP("rm").len)

	if mod == 0b11 {
		//var isS = isBitTrue(s)
		//var isWide = isBitTrue(w)
		var _, rm2 = transformRegToString(rm)
		var immediateValue = instructions[2]
		removeCount +=1
		fmt.Printf("%s %s %d \n", instructionType, rm2, immediateValue)
	}
	if mod == 0b00 {
		//var isS = isBitTrue(s)
		//var isWide = isBitTrue(w)
		var rm2 = translateMemoryDisplacement00(rm)
		var immediateValue = instructions[2]
		removeCount +=1
		if (reg == 0b000) {
			fmt.Printf("add [%s] %d \n", instructionType, rm2, immediateValue)
		}
		if (reg == 0b101) {
			fmt.Printf("sub [%s] %d \n", instructionType, rm2, immediateValue)
		}
		if (reg == 0b111) {
			fmt.Printf("cmp [%s] %d \n", instructionType, rm2, immediateValue)
		}
	}
	if mod == 0b10 {
		//var isS = isBitTrue(s)
		//var isWide = isBitTrue(w)
		var rm2 = translateMemoryDisplacement00(rm)
		var immediateValue = instructions[2]
		removeCount +=1
		if (reg == 0b000) {
			fmt.Printf("add [%s] %d \n", instructionType, rm2, immediateValue)
		}
		if (reg == 0b101) {
			fmt.Printf("sub [%s] %d \n", instructionType, rm2, immediateValue)
		}
		if (reg == 0b111) {
			fmt.Printf("cmp [%s] %d \n", instructionType, rm2, immediateValue)
		}
	}
	return instructions[removeCount:]
}

func disAddAccumulatorToRegister(instructionType string, instructions []uint8, inst inst) []uint8 {
	var removeCount = 0

	var instruction = instructions[0]
	removeCount += 1
	var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)

	if (isBitTrue(w)) {
		var uintResult uint16 = uint16(instructions[2]) << 8 | uint16(instructions[1]) 
		removeCount += 2
		fmt.Printf("%s ax %d \n", instructionType, uintResult)
	} else {
		removeCount += 1
		fmt.Printf("%s al %d \n", instructionType, instructions[1])
	}
	return instructions[removeCount:]
}

func clearBits(initial uint8, mask uint8, offset int) uint8 {
	value := initial & mask 
	return value >> offset
}

type inst struct {
	opType string
	opCode uint8
	opCodeDiff int
	params []P
}

type P struct{
	n string
	c uint8
	len int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func isBitTrue(input uint8) bool {
	if (input == 0b1) {
		return true
	}
	return false
}

func transformRegToString(reg uint8) (string, string) {
	var reg1 string
	var reg2 string
	switch (reg) {
		case 0b000:
			reg1 = "al"
			reg2 = "ax"
		case 0b001:
			reg1 = "cl"
			reg2 = "cx"
		case 0b010:
			reg1 = "dl"
			reg2 = "dx"
		case 0b011:
			reg1 = "bl"
			reg2 = "bx"
		case 0b100:
			reg1 = "ah"
			reg2 = "sp"
		case 0b101:
			reg1 = "ch"
			reg2 = "bp"
		case 0b110:
			reg1 = "dh"
			reg2 = "si"
		case 0b111:
			reg1 = "bh"
			reg2 = "di"
	}

	return reg1, reg2
}

func translateMemoryDisplacement00(reg uint8) string {
	var result string
	switch (reg) {
		case 0b000:
			result = "bx + si"
		case 0b001:
			result = "bx + di"
		case 0b010:
			result = "bp + si"
		case 0b011:
			result = "bp + di"
		case 0b100:
			result = "si"
		case 0b101:
			result = "di"
		case 0b110:
			result = "bp"
		case 0b111:
			result = "bx"
	}
	return result
}

