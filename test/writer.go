package main

import (
	hupwriter ".."
	"log"
	"time"
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
