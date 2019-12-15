package ks0108

//#cgo LDFLAGS: -lpigpio -lrt -pthread
//#include <pigpio.h>
//#include <stdint.h>
import "C"

import (
	"fmt"
	"math"
)

const (
	displaySetY = 0x40;
	displaySetX = 0xB8;
	displayStartLine = 0xC0;
	displayOnCmd = 0x3E;
)
	
// The Pins struct
type Pins struct {
	Rs uint8
	En uint8
	Cs1 uint8
	Cs2 uint8
	Cs3 uint8
	D0,D1,D2,D3,D4,D5,D6,D7 uint8
}

// The Ks0108 mains struct
type Ks0108 struct {
	pins Pins
	screenWidth, screenHeight uint8
	screenX, screenY uint8
	framebuffer []uint8
	fonts map[string][]uint8
}

// NewKs0108 initializes the screen and returns the instance
func NewKs0108(pins Pins, width uint8, height uint8) *Ks0108  {
	
	if (C.gpioInitialise() < 0) {
		fmt.Println("Ks0108 pigpio not initialized");
		return nil;
	}

	lcd := new(Ks0108);
	lcd.pins = pins;
	lcd.screenWidth = width;
	lcd.screenHeight = height;
	lcd.screenX, lcd.screenY = 0,0;
	lcd.framebuffer = make([]uint8, int(width) * int(height));
	lcd.fonts = make(map[string][]uint8);
	
	C.gpioSetMode(C.uint(pins.Rs), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.En), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.Cs1), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.Cs2), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.Cs3), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D0), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D1), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D2), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D3), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D4), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D5), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D6), C.PI_OUTPUT);
	C.gpioSetMode(C.uint(pins.D7), C.PI_OUTPUT);

	for i := uint8(0); i < 3; i++ {
		lcd.writeCommand((displayOnCmd | 0x01), i);
	}
	
	return lcd;
}

// ClearBuffer - clears the screen framebuffer
func (lcd *Ks0108) ClearBuffer() {
	for i := 0; i<len(lcd.framebuffer); i++ {
		lcd.framebuffer[i] = 0x00;
	}
}

// SyncBuffer - Sends data on framebuffer to the screen
func (lcd *Ks0108) SyncBuffer() {
	counter := 0;
	for row := uint8(0); row < 8; row++ {
		lcd.goTo(0,row);
		for col := uint8(0); col < lcd.screenWidth; col++ {
			lcd.writeData(lcd.framebuffer[counter]);
			counter++;
		}
	}
}

func (lcd *Ks0108) goTo(x uint8, y uint8) {
	var i uint8;
	lcd.screenX = x;
	lcd.screenY = y;
	for i=0; i<lcd.screenWidth/64; i++ {
		lcd.writeCommand(displaySetY | 0,i);
		lcd.writeCommand(displaySetX | y,i);
		lcd.writeCommand(displayStartLine | 0,i);
	}
	lcd.writeCommand(displaySetY | (x % 64), (x / 64));
	lcd.writeCommand(displaySetX | y, (x / 64));
}

func (lcd *Ks0108) putData(data uint8) {
	C.gpioWrite(C.uint(lcd.pins.D0), C.uint((data >> 0) & 1));
	C.gpioWrite(C.uint(lcd.pins.D1), C.uint((data >> 1) & 1));
	C.gpioWrite(C.uint(lcd.pins.D2), C.uint((data >> 2) & 1));
	C.gpioWrite(C.uint(lcd.pins.D3), C.uint((data >> 3) & 1));
	C.gpioWrite(C.uint(lcd.pins.D4), C.uint((data >> 4) & 1));
	C.gpioWrite(C.uint(lcd.pins.D5), C.uint((data >> 5) & 1));
	C.gpioWrite(C.uint(lcd.pins.D6), C.uint((data >> 6) & 1));
	C.gpioWrite(C.uint(lcd.pins.D7), C.uint((data >> 7) & 1));
}

func (lcd *Ks0108) writeData(dataToWrite uint8) {
	C.gpioWrite(C.uint(lcd.pins.Rs), 1);
	lcd.putData(dataToWrite);
	lcd.setController(lcd.screenX / 64, 1);
	C.gpioWrite(C.uint(lcd.pins.En), 1);
	C.gpioWrite(C.uint(lcd.pins.En), 0);
	lcd.setController(lcd.screenX / 64, 0);
	lcd.screenX++;
}

func (lcd *Ks0108) writeCommand(commandToWrite uint8, controller uint8) {
	C.gpioWrite(C.uint(lcd.pins.Rs), 0);
	lcd.setController(controller, 1);
	lcd.putData(commandToWrite);
	C.gpioWrite(C.uint(lcd.pins.En), 1);
	C.gpioWrite(C.uint(lcd.pins.En), 0);
	lcd.setController(controller, 0);
}

func (lcd *Ks0108) setController(controller uint8, enable uint8) {
	switch(controller){
		case 0 : C.gpioWrite(C.uint(lcd.pins.Cs1), C.uint(enable)); break;
		case 1 : C.gpioWrite(C.uint(lcd.pins.Cs2), C.uint(enable)); break;
		case 2 : C.gpioWrite(C.uint(lcd.pins.Cs3), C.uint(enable)); break;
	}
}

func (lcd *Ks0108) setPixel(x uint8, y uint8) {
	idx := int(float64(lcd.screenWidth) * math.Floor(float64(y)/8)) + int(x);
	lcd.framebuffer[idx] |= 1 << uint(y%8);
}

func (lcd *Ks0108) setPixels(x uint8, y uint8, data uint8) {
	idx := int(float64(lcd.screenWidth) * math.Floor(float64(y)/8)) + int(x);
	idx2 := int(float64(lcd.screenWidth) * (math.Floor(float64(y)/8) + 1)) + int(x);
	rest := uint8(y%8);
	lcd.framebuffer[idx] |=  data << uint(y%8);
	if(rest > 0) {
		lcd.framebuffer[idx2] |= data >> (uint(8) - uint(y%8));
	}
}