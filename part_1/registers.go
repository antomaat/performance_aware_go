package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	//data, err := os.ReadFile("./single_instruction")
	//check(err)
	data, err := ioutil.ReadFile("mov_single")
	Check(err)
	fmt.Printf("data is %0b\n", data)
	var index int = 0
	for index < len(data){
		var instruction uint8 = data[index]
		var register uint8 = data[index+1]
		fmt.Printf("load %08b and %08b\n", instruction, register)
		var opcode uint8 = ClearBits(instruction, 0b11111000) 
		switch (opcode) {
			case 0b10001000:
				DecodeMov(instruction, register)
			case 0b10110000:
				DecodeImmediateToRegister(instruction, register)
			case 0b10111000:
				DecodeImmediateToRegister(instruction, register)
				index +=1
		}
		index += 2
	}
}

const (
	source_is_reg uint8 = 0b00
	dest_is_reg uint8 = 0b10
)

const (
	is_8_bits uint8 = 0b0
	is_16_bits uint8 = 0b1
)

const(
	mod_mem_to_mem uint8 = 0b00
	mod_mem_8 uint8 = 0b01 //displacement mod on DISP-LO
	mod_mem_16 uint8 = 0b10 // displacement mod on DISP-LO and DISP-HI
	mod_reg_to_reg uint8 = 0b11 //register to register
)

func DecodeImmediateToRegister(instructions uint8, register uint8) {
	var w uint8 = ClearBits(instructions, 0b00001000) >> 3
	var reg uint8 = ClearBits(instructions, 0b00000111)
	var regString1, regString2 = TransformRegToString(reg)
	if (w == 0b00000001) {
		fmt.Printf("mov %s %d \n", regString2, register)
	}
	if (w == 0b00000000) {
		fmt.Printf("mov %s %d \n", regString1, register)
	}
}

func DecodeMov(instructions uint8, register uint8) {
	var direction uint8 = ClearBits(instructions, 0b00000010)
	var operation uint8 = ClearBits(instructions, 0b00000001)
	//var mode = ClearBits(register, 0b11000000) >> 6

	//LogDirection(direction)
	//LogOperation(operation)
	//LogModInfo(mode)


	var reg = ClearBits(register, 0b00111000)
	var reg1, reg2 = GetRegAsString(reg)

	var reverseReg = ClearBits(register, 0b00000111)
	var revReg1, revReg2 = GetRegAsString(reverseReg)

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

func LogModInfo(mod uint8) {
	switch mod {
	case mod_reg_to_reg:
		fmt.Println("mod is reg to reg")
	case mod_mem_8:
		fmt.Println("mod is 8 bits")
	}
}

func LogDirection(dir uint8) {
	switch dir {
	case source_is_reg:
		fmt.Println("source is reg")
	case dest_is_reg:
		fmt.Println("dest is reg")
	}
}

func LogOperation(dir uint8) {
	switch dir {
	case is_8_bits:
		fmt.Println("operation is 8 bits")
	case is_16_bits:
		fmt.Println("operation is 16 bits")
	}
}


func ClearBits(initial uint8, mask uint8) uint8 {
	return initial & mask
}

func TransformRegToString(reg uint8) (string, string) {
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

func GetRegAsString(reg uint8) (string, string) {
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

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

