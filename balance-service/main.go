package main

import (
	"balance-service/cmd/app"
	"log"
)

func main() {
	log.Println("startup...")
	err := app.Run()
	if err != nil {
		log.Fatal("error start application: ", err)
	}
	log.Println("finish")
	log.Println("adasda")
	log.Println("asdasd")
}
