package main

import (
	"log"
	"transfer-service/internal/bootstrap"
)

func main() {
	log.Println("startup...")
	err := bootstrap.Run()
	if err != nil {
		log.Fatal("error start application: ", err)
	}
}
