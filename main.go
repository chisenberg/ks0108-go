package main

import (
	"fmt"
	"ks0108-go/libs/ks0108"
	"time"
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

	lcd.LoadFont("metric01", "fonts/metric01.h");
	lcd.LoadFont("metric02", "fonts/metric02.h");
	lcd.LoadFont("metric04", "fonts/metric04.h");

	counter := 0;

	for {
		counter++;
		lcd.ClearBuffer();
		lcd.WriteString(64,6,"TESTANDO", "metric01");
		lcd.WriteString(0,0,"ABACATE", "metric02");
		lcd.WriteString(0,20, fmt.Sprintf("%d", counter), "metric04");
		lcd.SyncBuffer();
		time.Sleep(500 * time.Millisecond);
	}

}