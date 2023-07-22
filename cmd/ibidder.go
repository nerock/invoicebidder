package main

import (
	"fmt"
	"github.com/nerock/invoicebidder/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", cfg)
}
