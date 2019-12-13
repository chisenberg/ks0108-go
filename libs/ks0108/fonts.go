package ks0108

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"strings"
	"encoding/hex"
)

func loadFont(fileName string) []byte {
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
	
	return font;
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