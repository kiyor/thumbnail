package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kiyor/thumbnail"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	a := thumbnail.NewAgent()
	for _, v := range os.Args[1:] {
		err := filepath.Walk(v,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				ext := filepath.Ext(path)
				if ext == ".mp4" || ext == ".mkv" {
					dst := path + ".jpg"
					if _, err := os.Stat(dst); err == nil {
						log.Println("exist", dst)
						return nil
					}
					err = a.Thumbnail(path)
					if err != nil {
						return err
					}
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	}
}
