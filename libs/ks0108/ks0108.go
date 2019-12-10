package ks0108

//#cgo LDFLAGS: -lpigpio -lrt -pthread
//#include <pigpio.h>
//#include <stdint.h>
import "C"

import (
	"fmt"
	"time"
)

const (
	DISPLAY_SET_Y = 0x40;
	DISPLAY_SET_X = 0xB8;
	DISPLAY_START_LINE = 0xC0;
	DISPLAY_ON_CMD = 0x3E;
	ON = 0x01;
	OFF = 0x00;
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
	framebufferSize int
}

// InitKs0108 starts to listen
func InitKs0108(pins Pins, width uint8, height uint8) *Ks0108  {
	
	if (C.gpioInitialise() < 0) {
		fmt.Println("Ks0108 pigpio not initialized");
		return nil;
	}

	lcd := new(Ks0108);
	lcd.pins = pins;
	lcd.screenWidth = width;
	lcd.screenHeight = height;
	lcd.screenX, lcd.screenY = 0,0;
	lcd.framebufferSize = 1024;
	lcd.framebuffer = make([]uint8, lcd.framebufferSize);
	
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
		lcd.writeCommand((DISPLAY_ON_CMD | ON), i);
	}

	lcd.clearScreen();

	return lcd;
}

func (lcd *Ks0108) goTo(x uint8, y uint8) {
	var i uint8;
	lcd.screenX = x;
	lcd.screenY = y;
	for i=0; i<lcd.screenWidth/64; i++ {
		lcd.writeCommand(DISPLAY_SET_Y | 0,i);
		lcd.writeCommand(DISPLAY_SET_X | y,i);
		lcd.writeCommand(DISPLAY_START_LINE | 0,i);
	}
	lcd.writeCommand(DISPLAY_SET_Y | (x % 64), (x / 64));
	lcd.writeCommand(DISPLAY_SET_X | y, (x / 64));
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
	lcd.lcdDelay();
	C.gpioWrite(C.uint(lcd.pins.Rs), 1);
	lcd.putData(dataToWrite);
	lcd.setController(lcd.screenX / 64, 1);
	C.gpioWrite(C.uint(lcd.pins.En), 1);
	lcd.lcdDelay();
	C.gpioWrite(C.uint(lcd.pins.En), 0);
	lcd.setController(lcd.screenX / 64, 0);
	lcd.screenX++;
}

func (lcd *Ks0108) writeCommand(commandToWrite uint8, controller uint8) {
	lcd.lcdDelay();
	C.gpioWrite(C.uint(lcd.pins.Rs), 0);
	lcd.setController(controller, 1);
	lcd.putData(commandToWrite);
	C.gpioWrite(C.uint(lcd.pins.En), 1);
	lcd.lcdDelay();
	C.gpioWrite(C.uint(lcd.pins.En), 0);
	lcd.setController(controller, 0);
}

func (lcd *Ks0108) lcdDelay() {
	time.Sleep(5 * time.Microsecond)
}

func (lcd *Ks0108) setController(controller uint8, enable uint8) {
	switch(controller){
		case 0 : C.gpioWrite(C.uint(lcd.pins.Cs1), C.uint(enable)); break;
		case 1 : C.gpioWrite(C.uint(lcd.pins.Cs2), C.uint(enable)); break;
		case 2 : C.gpioWrite(C.uint(lcd.pins.Cs3), C.uint(enable)); break;
	}
}

func (lcd *Ks0108) clearScreen() {
	lcd.clearBuffer();
	lcd.syncBuffer();
}

func (lcd *Ks0108) clearBuffer() {
	for i := 0; i<lcd.framebufferSize; i++ {
		lcd.framebuffer[i] = 0x00;
	}
}

func (lcd *Ks0108) syncBuffer() {
	counter := 0;
	fmt.Println(len(lcd.framebuffer));
	for row := uint8(0); row < 8; row++ {
		lcd.goTo(0,row);
		for col := uint8(0); col < lcd.screenWidth; col++ {
			lcd.writeData(lcd.framebuffer[counter]);
			counter++;
		}
	}
}

// //-------------------------------------------------------------------------------------------------
// //
// //-------------------------------------------------------------------------------------------------
// void Ks0108pi::setPixel(uint8_t x, uint8_t y)
// {
// 	int idx = (SCREEN_WIDTH * (y/8)) + x;
// 	framebuffer[idx] |= 1 << y%8;
// }

// //-------------------------------------------------------------------------------------------------
// //
// //-------------------------------------------------------------------------------------------------
// void Ks0108pi::clearPixel(uint8_t x, uint8_t y)
// {
// 	int idx = (SCREEN_WIDTH * (y/8)) + x;
// 	framebuffer[idx] &= ~(1 << y%8);
// }

// //-------------------------------------------------------------------------------------------------
// //
// //-------------------------------------------------------------------------------------------------
// void Ks0108pi::setPixels(uint8_t x, uint8_t y, uint8_t byte)
// {
// 	int idx = (SCREEN_WIDTH * (y/8)) + x;
// 	int idx2 = (SCREEN_WIDTH * ( (y/8)+1) ) + x;
// 	uint8_t rest = y%8;
// 	framebuffer[idx] |= ( byte << y%8 );
// 	if(rest)
// 		framebuffer[idx2] |= byte >> (8-y%8);
// }

// //-------------------------------------------------------------------------------------------------
// //
// //-------------------------------------------------------------------------------------------------
// void Ks0108pi::drawRect(uint8_t x, uint8_t y, uint8_t w, uint8_t h, uint8_t style){
// 	for(int nx=x; nx < (x+w) ; nx++){
// 		for(int ny=y; ny < (y+h) ; ny++){
// 			if(style & STYLE_BLACK_BG) setPixel(nx,ny);
// 			else if(style & STYLE_WHITE_BG) clearPixel(nx,ny);
// 			else if(style & STYLE_GRAY_BG){
// 				if((nx+ny)%2==0)
// 					clearPixel(nx,ny);
// 				else
// 					setPixel(nx,ny);
// 			}
// 		}
// 	}

// 	if( (style & STYLE_BLACK_BORDER) || (style & STYLE_WHITE_BORDER)) {
// 		drawLine(x,y,(x+w)-1,y);
// 		drawLine(x,(y+h)-1,(x+w)-1,(y+h)-1);
// 		drawLine(x,y,x,(y+h)-1);
// 		drawLine((x+w)-1,y,(x+w)-1,(y+h)-1);
// 	}
// }


// //-------------------------------------------------------------------------------------------------
// //
// //-------------------------------------------------------------------------------------------------
// void Ks0108pi::drawLine(uint8_t x0, uint8_t y0, uint8_t x1, uint8_t y1){
// 	int dx = abs(x1-x0), sx = x0<x1 ? 1 : -1;
// 	int dy = abs(y1-y0), sy = y0<y1 ? 1 : -1;
// 	int err = (dx>dy ? dx : -dy)/2, e2;

// 	for(;;){
// 		setPixel(x0,y0);
// 		if (x0==x1 && y0==y1) break;
// 		e2 = err;
// 		if (e2 >-dx) { err -= dy; x0 += sx; }
// 		if (e2 < dy) { err += dx; y0 += sy; }
// 	}

// }

// //-------------------------------------------------------------------------------------------------
// //
// //-------------------------------------------------------------------------------------------------
// void Ks0108pi::writeChar(uint8_t x, uint8_t y, char charToWrite, uint8_t* font)
// {
// 	int firstChar = font[4];
// 	int charCount = font[5];
// 	int charHeight = font[3];
// 	int charWidth = font[2];
// 	int sum= 6;
// 	int fixed_width = 1;

// 	if( (font[0] + font [1]) != 0x00){
// 		fixed_width  = 0;
// 	}


// 	if( !fixed_width ){
// 		charWidth = font[6+(charToWrite-firstChar)];
// 		sum += charCount;
// 	}

// 	//jumps to the char data position on the array.
// 	for(int i=firstChar; i<charToWrite; i++){
// 		if( !fixed_width )
// 			sum += font[6+i-firstChar] * ceil(charHeight/8.0);
// 		else
// 			sum += charWidth * ceil(charHeight/8.0);
// 	}

// 	for(int line=0; line < charHeight; line+=8){
// 		for(int col=0; col<charWidth; col++){
// 			setPixels(x+col, ceil(y+line),
// 				font[sum + col + (int)ceil(charWidth*line/8.0)]
// 			);
// 		}
// 	}

// }



// //-------------------------------------------------------------------------------------------------
// //
// //-------------------------------------------------------------------------------------------------
// void Ks0108pi::writeString(uint8_t x, uint8_t y, char * stringToWrite, uint8_t* font)
// {
// 	while(*stringToWrite){
// 		writeChar(x,y,*stringToWrite++, font);
// 		x+=font[2]+1;
// 	}
// }