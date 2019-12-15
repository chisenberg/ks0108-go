package main

import (
	// "fmt"
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
		// rain
		lcd.WriteChar(0,0, byte(0),"icons");
		lcd.WriteString(14,0,"23.5", "metric02");

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

		// lcd.WriteString(85,0,"TESTANDO", "metric01");
		// lcd.WriteString(85,8,"KS0108", "metric01");
		// lcd.WriteString(0,20, fmt.Sprintf("%d", counter), "metric04");
		// // lcd.WriteString(0,20, "11:38", "metric04");
		// lcd.DrawLine(0,16,127,16);
		// lcd.DrawRect(64,32,30,30,false);
		
		lcd.SyncBuffer();
		time.Sleep(500 * time.Millisecond);
	}

}
