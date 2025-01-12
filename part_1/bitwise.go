package main

import "fmt"

func main() {
	//var x uint8 = 0b00000
	//fmt.Printf("%08b\n", x)
	//fmt.Println(x)
	shifting()
	//lightManipulation()
}

func shifting() {
	var original uint8 = 0b11010111
	var shift_left uint8 = original << 6
	var shift_right uint8 = original >> 6 & trunc(2)
	var isOne = shift_right == 3
	fmt.Printf("original %03b\n", original)
	fmt.Printf("shift left %08b\n", shift_left)
	fmt.Printf("shift right %08b\n", shift_right)
	fmt.Printf("isOne\n", isOne)
}

func lightManipulation() {
	var lights uint8 = 0b10001001
	lights = turnAllLightsOff(lights)
	//lights = toggleOutsideLights(lights)
	//lights = turnOnOutsideLights(lights)
	fmt.Printf("lights turnt on %08b\n", lights)

	// use the and operation to check with the mask if the selected ligts are on
	var isLightsOn = isLightsOn(lights)

	if (isLightsOn) {
		fmt.Println("At least one outside ligtht is on")
	} else {
		fmt.Println("All lights are off")
	}
	fmt.Printf("%08b\n", lights)

}

func isLightsOn(lights uint8) bool {
	var mask uint8 = 0b00110000
	var isLightsOn = lights & mask
	return isLightsOn != 0
}

func turnOffOutsideLight(lights uint8) uint8 {
	var mask uint8 = 0b00110000
	var invertMask = ^mask
	return  lights & invertMask
}

func turnOnOutsideLights(lights uint8) uint8 {
	var mask uint8 = 0b00110000
	return lights | mask
}

func toggleOutsideLights(lights uint8) uint8 {
	var mask uint8 = 0b00110000
	return lights ^ mask
}

func turnAllLightsOff(lights uint8) uint8 {
	return lights ^ lights
}

func trunc(length int) byte {
	m := byte(1)
	for i := 1; i < length; i++ {
		m = m << 1
		m |= 1
	}
	return m
}
