package main

import (
	"flag"
	"log"
)

const (
	fetcherCount      = 50
	parserCount       = 2
	inputQueueSize    = 0
	fetchQueueSize    = 0
	fetchedQueueSize  = 0
	storageQueueSize  = 0
	rawQueueSize      = 0
	filteredQueueSize = 0
)

func main() {
	flag.Parse()
	noServer := flag.Bool("noserver", false, "if set, no local server will be started")
	publicServer := flag.Bool("publicserver", false, "if set, the server will listen on all interfaces")

	inputQueue, fetchQueue := setupURLQueue()
	fetchedQueue := setupFetchers(fetcherCount, fetchQueue)
	storageQueue := setupStorage(fetchedQueue)
	setupParsers(parserCount, inputQueue, storageQueue)
	startupServer(*noServer, *publicServer, inputQueue)

	log.Println("Main thread done...")
}
