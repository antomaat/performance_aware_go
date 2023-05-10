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
			P{c: 0b00000010, n: "s", len: 1},
			P{c: 0b00000001, n: "w", len: 0},
			P{c: 0b11000000, n: "mod", len: 6},
			P{c: 0b00111000, n: "reg", len: 3},
			P{c: 0b00000111, n: "rm", len: 0},
		},
	},
} 

func (inst inst) getP(value string) P {
	for index := 0; index < len(inst.params); index++ {
		if inst.params[index].n == value {
			return inst.params[index]
		}
	}
	panic("no such param for inst")
}

func main() {
	//data, err := os.ReadFile("./single_instruction")
	//check(err)
	//data, err := ioutil.ReadFile("single_instruction")
	//data, err := ioutil.ReadFile("multiple_instructions")
	//data, err := ioutil.ReadFile("mov_inst_complex")
	data, err := ioutil.ReadFile("add_sub_cmp")
	fmt.Printf("inst %08b\n", data)
	check(err)
	var max = 30
	//for len(data) > 0 && max > 0 {
	for len(data) > 0 {
		data = disassembleAndReturn(data)
		max -= 1
		//fmt.Printf("inst %08b\n", data)
	}
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
				return disAddComplex(instructions, inst)
			}
			if inst.opType == "addImmReg" {
				return disAddImmediateToRegister(instructions, inst)
			}
			if inst.opType == "addImmAcc" {
				return disAddAccumulatorToRegister(instructions, inst)
			}
		}
	}
	return instructions
}

//only one byte big operation
func disMovImmediateToRegister(instructions []uint8, inst inst) []uint8 {
	var removeCount = 0
	// first one is for opCode
	var instruction = instructions[0]
	var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)
	var reg = clearBits(instruction, inst.getP("reg").c, inst.getP("reg").len)

	//fmt.Printf("w %08b\n", w)
	var short, wide = transformRegToString(reg)

	if (isBitTrue(w)){
		removeCount += 2
		var result uint16 = uint16(instructions[2]) << 8 | uint16(instructions[1]) 
		//fmt.Printf("result %016b\n", result)
		fmt.Printf("mov %s %d\n", wide, result)
	} else {
		removeCount += 1
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

func disAddComplex(instructions []uint8, inst inst) []uint8 {
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
		fmt.Printf("add %s %s \n", destination, source)
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
		fmt.Printf("add %s, %s \n", destination, source)
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
		fmt.Printf("add %s, %s \n", destination, source )
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
		fmt.Printf("add %s %s \n", destination, source)
		removeCount += 2
	}

	return instructions[removeCount:]
}

func disAddImmediateToRegister(instructions []uint8, inst inst) []uint8 {
	var removeCount = 0
	// first one is for opCode
	//var instruction = instructions[0]
	//var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)
	//var s = clearBits(instruction, inst.getP("s").c, inst.getP("s").len)
	removeCount +=1

	var registers = instructions[1]
	removeCount += 1
	var mod = clearBits(registers, inst.getP("mod").c, inst.getP("mod").len)
	//var reg = clearBits(registers, inst.getP("reg").c, inst.getP("reg").len)
	var rm = clearBits(registers, inst.getP("rm").c, inst.getP("rm").len)

	if mod == 0b11 {
		//var isS = isBitTrue(s)
		//var isWide = isBitTrue(w)
		var _, rm2 = transformRegToString(rm)
		var immediateValue = instructions[2]
		removeCount +=1
		fmt.Printf("add %s %d \n", rm2, immediateValue)
	}
	if mod == 0b00 {
		//var isS = isBitTrue(s)
		//var isWide = isBitTrue(w)
		var rm2 = translateMemoryDisplacement00(rm)
		var immediateValue = instructions[2]
		removeCount +=1
		fmt.Printf("add [%s] %d \n", rm2, immediateValue)
	}
	if mod == 0b10 {
		//var isS = isBitTrue(s)
		//var isWide = isBitTrue(w)
		var rm2 = translateMemoryDisplacement00(rm)
		var immediateValue = instructions[2]
		removeCount +=1
		fmt.Printf("add [%s] %d \n", rm2, immediateValue)
	}
	return instructions[removeCount:]
}

func disAddAccumulatorToRegister(instructions []uint8, inst inst) []uint8 {
	var removeCount = 0

	var instruction = instructions[0]
	removeCount += 1
	var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)

	if (isBitTrue(w)) {
		var uintResult uint16 = uint16(instructions[2]) << 8 | uint16(instructions[1]) 
		removeCount += 2
		fmt.Printf("add ax %d \n", uintResult)
	} else {
		removeCount += 1
		fmt.Printf("add al %d \n", instructions[1])
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

