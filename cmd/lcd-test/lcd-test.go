package main

import (
	"fmt"
	"time"
	"github.com/chisenberg/ks0108-go/ks0108"
)

func main() {

	conf := ks0108.Pins{
		Rs: 7,
		En: 11,
		Cs1: 25,
		Cs2: 8,
		Cs3: 9,
		D0: 2,
		D1: 3,
		D2: 4,
		D3: 17,
		D4: 27,
		D5: 22,
		D6: 10,
		D7: 9,
	}

	lcd := ks0108.NewKs0108(conf, 128, 64);

	lcd.LoadFont("metric01", "assets/fonts/metric01.h");
	lcd.LoadFont("metric02", "assets/fonts/metric02.h");
	lcd.LoadFont("metric04", "assets/fonts/metric04.h");
	lcd.LoadFont("icons", "assets/fonts/icons12.h");

	counter := 0;

	for {
		counter++;
		lcd.ClearBuffer();

		lcd.WriteString(66,0,"MOWS", "metric02");
		lcd.WriteString(100,0,"STATION", "metric01");
		lcd.WriteString(108,7,"22:53", "metric01");

		// rain
		lcd.WriteChar(0,0, byte(0),"icons");
		lcd.WriteString(14,0,"0.0", "metric02");

		// temp
		lcd.WriteChar(0,18, byte(1),"icons");
		lcd.WriteString(14,16,"38.2", "metric02");
		lcd.WriteString(64,16,"MIN 28.1'", "metric01");
		lcd.WriteString(64,23,"MAX 56.4'", "metric01");

		// humid
		lcd.WriteChar(0,32, byte(2),"icons");
		lcd.WriteString(14,32,"38.2%", "metric02");
		lcd.WriteString(64,32,"MIN 28.1%", "metric01");
		lcd.WriteString(64,39,"MAX 56.4%", "metric01");

		// wind
		lcd.WriteChar(0,48, byte(3),"icons");
		lcd.WriteString(14,48,"38.2", "metric02");
		lcd.WriteString(64,48,"MIN ----", "metric01");
		lcd.WriteString(64,55,"MAX 56.4", "metric01");

		lcd.DrawLine(64,0,64,13);
		lcd.DrawLine(65,13,104,13);
		lcd.DrawLine(105,13,105,63);

		lcd.WriteChar(112,15, byte(4),"icons");
		lcd.WriteString(108,29,fmt.Sprintf("%d",counter), "metric01");

		lcd.DrawLine(105,37,127,37);

		lcd.WriteChar(112,40, byte(5),"icons");
		lcd.WriteString(108,55,"4.52V", "metric01");
		
		lcd.SyncBuffer();
		time.Sleep(500 * time.Millisecond);
	}

}
