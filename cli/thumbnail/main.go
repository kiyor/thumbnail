package main

import (
	"log"
	"os"

	"github.com/kiyor/thumbnail"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	a := thumbnail.NewAgent()
	for _, v := range os.Args[1:] {
		err := a.Thumbnail(v)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
}
