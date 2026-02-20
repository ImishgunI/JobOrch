package main

import (
	"log"
	"time"
)

const maxConcurrency = 3

func main() {
	sem := make(chan struct{}, maxConcurrency)
	for i := range 5 {
		sem <- struct{}{}
		log.Printf("Start job %d\n", i)
		go func(i int) {
			defer func(i int) {
				log.Printf("Finish Job %d\n", i)
				<-sem
			}(i)
			time.Sleep(2 * time.Second)
		}(i)
	}
}
