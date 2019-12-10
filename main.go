package main

import (
	// "fmt"
	"ks0108-go/libs/ks0108"
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

	ks0108.InitKs0108(conf, 128, 64);

}