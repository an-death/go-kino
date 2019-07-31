package main

import (
	"log"
	"os"

	"github.com/an-death/go-kino/frontend"
)

var PORT string

func init() {
	PORT = os.Getenv("PORT")

	if PORT == "" {
		PORT = "8000"
	}
}

func main() {
	log.Fatal(frontend.Run(":" + PORT))
}
