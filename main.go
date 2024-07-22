// main.go

package main

import (
	"log"
	"os"

	apifolder "github.com/gohugoio/hugo/apifolder"
	"github.com/gohugoio/hugo/commands"
)

func main() {
	log.SetFlags(0)
	go apifolder.StartServer() // Start the API server in a goroutine
	err := commands.Execute(os.Args[1:])
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
