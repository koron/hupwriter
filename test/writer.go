package main

import (
	"log"
	"time"

	"github.com/koron/hupwriter"
)

func main() {
	w, err := hupwriter.New("myapp.log", "myapp.pid")
	if err != nil {
		panic("can't open log file: " + err.Error())
	}
	log.SetOutput(w)
	count := 0
	for {
		log.Printf("%d", count)
		count++
		time.Sleep(time.Second)
	}
}
