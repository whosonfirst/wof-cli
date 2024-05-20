package main

import (
	"context"
	"os"
	"fmt"
	"log"
	
	"github.com/whosonfirst/wof/app/format"
)

func usage(){
	fmt.Println("Usage")
	os.Exit(0)
}

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	ctx := context.Background()
	
	cmd := os.Args[1]

	// To do: Need to make os.Args[2:] the flag/cli input to 'cmd'
	
	var err error
	
	switch cmd {
	case "fmt":
		err = format.Run(ctx)
	default:
		usage()
	}

	if err != nil {
		log.Fatalf("Failed to run '%s' command, %w", cmd, err)
	}

	os.Exit(0)
}
