package main

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Config: Error loading .env file")
	}
}
func main() {

}
