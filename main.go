package main

import (
	"fmt"
	"io/ioutil"
)

var instTable = []inst{
	inst{
		opType: "mov", 
		opCode: 0b100010,
		params: []P{
			P{c: 0b00000010, n: "d", len: 1},
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
	data, err := ioutil.ReadFile("multiple_instructions")
	check(err)
	var index int = 0
	for index < len(data){
	//for index < 4 { 
		var instruction uint8 = data[index]
		var register uint8 = data[index+1]
		disassemble(instruction, register)
		index +=2
	}
	for len(data) > 0 {
		data = disassembleAndReturn(data)
	}
}

func disassembleAndReturn(instructions []uint8) []uint8 {
	var opCode = instructions[0] >> 2
	if (opCode == instTable[0].opCode) {
		return disMovComplex(instructions, instTable[0])
	}
	return instructions
}

func disassemble(instruction uint8, register uint8) {
	var opCode = instruction >> 2
	if (opCode == instTable[0].opCode) {
		disMov(instruction, register, instTable[0])
	}
}

func disMovComplex(instructions []uint8, inst inst) []uint8 {
	var removeCount = 0
	// first one is for opCode
	var instruction = instructions[0]
	//every move should have a register part?
	var register = instructions[1]


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

	removeCount += 1
	return instructions[removeCount:]
}

func disMov(instruction uint8, register uint8, inst inst) {
	var d = clearBits(instruction, inst.getP("d").c, inst.getP("d").len)
	var w = clearBits(instruction, inst.getP("w").c, inst.getP("w").len)
	var mod = clearBits(register, inst.getP("mod").c, inst.getP("mod").len)
	var reg = clearBits(register, inst.getP("reg").c, inst.getP("reg").len)
	var rm = clearBits(register, inst.getP("rm").c, inst.getP("rm").len)
	//var reg = inst.getP("reg")
	//fmt.Printf("instruction %08b \n", instruction)
	//fmt.Printf("register %08b \n", register)
	//fmt.Printf("d %08b \n", d)
	//fmt.Printf("w %08b \n", w)
	//fmt.Printf("mod %08b \n", mod)
	//fmt.Printf("reg %08b \n", reg)
	//fmt.Printf("rm %08b \n", rm)
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
}

func clearBits(initial uint8, mask uint8, offset int) uint8 {
	value := initial & mask 
	return value >> offset
}

type inst struct {
	opType string
	opCode uint8
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

