package main

import (
	"log"
	"os"

<<<<<<< HEAD
	"github.com/an-death/go-kino/frontend"
=======
	"github.com/an-death/go-kino/providers/releases"
>>>>>>> dev
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
