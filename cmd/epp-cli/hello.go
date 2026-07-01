package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runHello(
	client *epp.Client,
	enabled bool,
) error {

	if !enabled {
		return nil
	}

	greeting, err := client.Hello()
	if err != nil {
		return err
	}

	printGreeting("HELLO", greeting)

	return nil
}

func printGreeting(
	title string,
	greeting *types.Greeting,
) {

	fmt.Println()
	fmt.Printf("========== %s ==========\n", title)

	fmt.Printf("Server ID            : %s\n", greeting.ServerID)
	if greeting.ServerDate != nil {
		fmt.Printf("Server Date          : %s\n", greeting.ServerDate.Format(time.RFC3339))
	}
	fmt.Printf("Supported Objects    : %s\n", strings.Join(greeting.SupportedObjects, ", "))
	fmt.Printf("Supported Extensions : %s\n", strings.Join(greeting.SupportedExtensions, ", "))
	fmt.Printf("Languages            : %s\n", strings.Join(greeting.Languages, ", "))
	fmt.Printf("Versions             : %s\n", strings.Join(greeting.Versions, ", "))
}
