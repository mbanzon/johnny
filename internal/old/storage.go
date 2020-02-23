package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type storedFetchedURL struct {
	URL         string
	Data        []byte
	ContentType string
}

func setupStorage(fetchedQueue chan *fetchedURL) (storageQueue chan *fetchedURL) {
	storageQueue = make(chan *fetchedURL, storageQueueSize)

	var current *fetchedURL
	ticker := time.NewTicker(100 * time.Millisecond)
	outID := 0
	inID := 0

	go func() {
		for {
			select {
			case tmp := <-fetchedQueue:
				if current == nil && outID == inID {
					current = tmp
				} else {
					fp, err := os.Create(fmt.Sprintf("stored_%010d.json", outID))
					if err != nil {
						log.Println("error creating storage file:", err)
						break
					}
					encoder := json.NewEncoder(fp)
					err = encoder.Encode(storedFetchedURL{URL: tmp.url, Data: tmp.data.Bytes(), ContentType: tmp.contentType})
					if err != nil {
						log.Println("error encoding data to storage file:", err)
						fp.Close()
						break
					}
					fp.Close()
					outID++
				}
				break
			case <-ticker.C:
				break
			}

			if current == nil && outID > inID {
				fp, err := os.Open(fmt.Sprintf("stored_%010d.json", inID))
				if err != nil {
					log.Println("error opening storage file:", err)
				} else {
					decoder := json.NewDecoder(fp)
					var tmp storedFetchedURL
					err := decoder.Decode(&tmp)
					if err != nil {
						log.Println("error decoding data from storage file:", err)
						fp.Close()
					} else {
						current = &fetchedURL{url: tmp.URL, data: bytes.NewBuffer(tmp.Data), contentType: tmp.ContentType}
						fp.Close()
						err := os.Remove(fmt.Sprintf("stored_%010d.json", inID))
						if err != nil {
							panic(err)
						}
						inID++
					}
				}
			}

			if current != nil {
				select {
				case storageQueue <- current:
					current = nil
					break
				case <-ticker.C:
					break
				}
			}
		}
	}()

	return storageQueue
}
