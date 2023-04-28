package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	//data, err := os.ReadFile("./single_instruction")
	//check(err)
	data, err := ioutil.ReadFile("multiple_instructions")
	check(err)
	var index int = 0
	for index < len(data){
		var instruction uint8 = data[index]
		var register uint8 = data[index+1]
		var opcode uint8 = clearBits(instruction, 0b11111100) 
		if (opcode == 0b10001000) {
			decodeMov(instruction, register)
		}
		index += 2
	}
}

func decodeMov(instructions uint8, register uint8) {
	var direction uint8 = clearBits(instructions, 0b00000010)
	var operation uint8 = clearBits(instructions, 0b00000001)

	var mode = clearBits(register, 0b11000000)
	var memoryMode = isMemoryMode(mode)
	if (memoryMode) {
		fmt.Println("is from memory\n")
	}

	var reg = clearBits(register, 0b00111000)
	var reg1, reg2 = getRegAsString(reg)

	var reverseReg = clearBits(register, 0b00000111)
	var revReg1, revReg2 = getRegAsString(reverseReg)

	if (operation == 0b00000001) {
		if (direction == 0b00000010) {
			fmt.Printf("mov %s %s\n", reg2, revReg2)
		}
		if (direction == 0b00000000) {
			fmt.Printf("mov %s %s\n", revReg2, reg2)
		}
	} else {
		if (direction == 0b00000010) {
			fmt.Printf("mov %s %s\n", reg1, revReg1)
		}
		if (direction == 0b00000000) {
			fmt.Printf("mov %s %s\n", revReg1, reg1)
		}
	}
}

func isMemoryMode(input uint8) bool {
	if (input == 0b00000000) {
		return true
	}
	return false
}

func clearBits(initial uint8, mask uint8) uint8 {
	return initial & mask
}

func getRegAsString(reg uint8) (string, string) {
	var reg1 string
	var reg2 string
	switch reg {
	case 0b00000000:
		reg1 = "al"
		reg2 = "ax"
	case 0b00001000, 0b00000001:
		reg1 = "cl"
		reg2 = "cx"
	case 0b00010000, 0b00000010:
		reg1 = "dl"
		reg2 = "dx"
	case 0b00011000, 0b00000011:
		reg1 = "bl"
		reg2 = "bx"
	case 0b00100000, 0b00000100:
		reg1 = "ah"
		reg2 = "sp"
	case 0b00101000, 0b00000101:
		reg1 = "ch"
		reg2 = "bp"
	case 0b00110000, 0b00000110:
		reg1 = "dh"
		reg2 = "si"
	case 0b00111000, 0b00000111:
		reg1 = "bh"
		reg2 = "di"
	}
	return reg1, reg2
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

