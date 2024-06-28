package main

import (
	"log"
	"time"

	"github.com/koron/hupwriter"
)

func main() {
	w := hupwriter.New("myapp.log", "myapp.pid")
	log.SetOutput(w)
	count := 0
	for {
		log.Printf("%d", count)
		count++
		time.Sleep(time.Second)
	}
}
