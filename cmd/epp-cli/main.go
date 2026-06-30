package main

import (
	"fmt"
	"log"

	"github.com/hariom-pal/go-epp/internal/config"
	"github.com/hariom-pal/go-epp/internal/epp"
)

func main() {

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("TLS Connected Successfully")

	fmt.Println("======================================")
	fmt.Println(string(client.Greeting()))
	fmt.Println("======================================")

	err = client.Login()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Logout()
	if err != nil {
		log.Fatal(err)
	}
}
