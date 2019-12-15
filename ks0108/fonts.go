package ks0108

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"strings"
	"encoding/hex"
	"math"
)

// LoadFont - loads font from file into Ks0108 instance
func (lcd *Ks0108) LoadFont(fontName string, fileName string) {
	
	// abort if font is loaded
	if _, ok := lcd.fonts[fontName]; ok {
		return;
	}

	var font []byte;
	file, err := os.Open(fileName);
	if err != nil {
		log.Fatal(err);
	}
	defer file.Close();

	scanner := bufio.NewScanner(file);
	// split by commas and lines
	scanner.Split(onComma);

	var token string;
	for scanner.Scan() {
		token = strings.TrimSpace(scanner.Text());
		// ignore comments and white lines
		if strings.HasPrefix(token, "//") || strings.HasPrefix(token, "#") || len(token) < 1 {
			continue;
		}
		// remove 0x from byte
		token := strings.Replace(token, "0x", "", -1);
		// add 0 in front if length is 1
		if len(token) == 1 {
			token = "0" + token;
		}
		// decode
		decoded, err := hex.DecodeString(token)
		if err != nil {
			fmt.Println(token)
			log.Fatal(err)
		}
		font = append(font, decoded[0]);
	}

	lcd.fonts[fontName] = font;
}

// WriteString - write string into x,y position
func (lcd *Ks0108) WriteString(x uint8, y uint8, stringToWrite string, fontName string) {
	font := lcd.fonts[fontName];
	for idx := range stringToWrite {
		lcd.WriteChar(x, y, stringToWrite[idx], fontName);
		x+=font[2]+1;
	}
}

// WriteChar - write single char or icon into x,y position
func (lcd *Ks0108) WriteChar(x uint8, y uint8, charToWrite byte, fontName string) {
	font := lcd.fonts[fontName];
	firstChar := font[4];
	charCount := int(font[5]);
	charHeight := font[3];
	charWidth := font[2];
	sum := int(6); // 6 bytes of header
	fixedWidth := true;

	if( (font[0] + font [1]) != 0x00){
		fixedWidth  = false;
	}

	if( !fixedWidth ){
		charWidth = font[6+(charToWrite-firstChar)];
		sum += charCount;
	}

	//jumps to the char data position on the array.
	for i:=firstChar; i<charToWrite; i++ {
		if( !fixedWidth ) {
			sum += int(float64(font[6+i-firstChar]) * math.Ceil(float64(charHeight)/8.0));
		} else {
			sum += int(float64(charWidth) * math.Ceil(float64(charHeight)/8.0));
		}
	}

	for line:=uint8(0); line < uint8(charHeight); line+=8 {
		for col:=uint8(0); col< uint8(charWidth); col++ {
			setByte := font[sum + int(col) + int(math.Ceil(float64(charWidth)*float64(line/8.0)))];
			lcd.setPixels(x+col, y+line, uint8(setByte));
		}
	}

}

func onComma(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == ',' || data[i] == '\n' {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}