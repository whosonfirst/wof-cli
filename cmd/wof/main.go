package main

import (
	"context"
	"log"
	"os"

	"github.com/whosonfirst/wof/app/wof"
)

func main() {

	ctx := context.Background()
	err := wof.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
