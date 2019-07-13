package main

import (
	"log"
	"time"
)

func setupURLQueue() (chan []string, chan string) {
	log.Println("Setting up fetch queue...")

	rawQueue, filteredQueue := setupURLQueueFilter()

	var queue []string
	fetchQueue := make(chan string, fetchQueueSize)

	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			select {
			case u := <-filteredQueue:
				queue = append(queue, u...)
				break
			case <-ticker.C:
				break
			}

			if len(queue) > 0 {
				select {
				case fetchQueue <- queue[0]:
					queue = queue[1:]
					break
				case <-ticker.C:
					break
				}
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("Queue length:", len(queue))
		}
	}()

	return rawQueue, fetchQueue
}

func setupURLQueueFilter() (rawQueue chan []string, filteredQueue chan []string) {
	rawQueue = make(chan []string, rawQueueSize)
	filteredQueue = make(chan []string, filteredQueueSize)

	lookup := make(map[string]bool)

	go func() {
		for {
			var filtered []string
			links := <-rawQueue
			for _, l := range links {
				if _, found := lookup[l]; !found {
					filtered = append(filtered, l)
					lookup[l] = true
				}
			}

			log.Printf("Got %d new links (%d dups)", len(filtered), len(links)-len(filtered))
			filteredQueue <- filtered
		}
	}()

	return rawQueue, filteredQueue
}
