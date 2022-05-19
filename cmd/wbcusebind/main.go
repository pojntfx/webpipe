package main

import (
	"log"
	"os"

	"github.com/pojntfx/webpipe/pkg/cuse"
)

func main() {
	log.Println(
		cuse.Args{
			Argc: int32(len(os.Args)),
		},
	)
}
