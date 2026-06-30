package main

import (
	"fmt"

	"github.com/hariom-pal/go-epp/pkg/idn"
)

func testIDN() {

	tests := []string{
		"google.in",
		"भारत.in",
		"भारत.भारत",
		"ગુજરાત.ભારત",
		"தமிழ்.இந்தியா",
		"বাংলা.ভারত",
		"ਪੰਜਾਬੀ.ਭਾਰਤ",
		"ଓଡ଼ିଆ.ଭାରତ",
	}

	for _, d := range tests {

		ascii, err := idn.ToASCII(d)
		if err != nil {
			fmt.Println(err)
			continue
		}

		unicode, _ := idn.ToUnicode(ascii)

		fmt.Println("--------------------------------")
		fmt.Println("Input    :", d)
		fmt.Println("ASCII    :", ascii)
		fmt.Println("Unicode  :", unicode)
	}
}
