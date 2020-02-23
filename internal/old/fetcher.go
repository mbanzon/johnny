package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
)

type fetchedURL struct {
	url         string
	data        *bytes.Buffer
	contentType string
}

func setupFetchers(count int, fetchQueue chan string) (fetchedQueue chan *fetchedURL) {
	log.Println("Setting up fetchers:", count)

	fetchedQueue = make(chan *fetchedURL, fetchedQueueSize)

	for i := 0; i < count; i++ {
		go func() {
			for {
				u := <-fetchQueue

				res, err := http.Get(u)
				if err != nil {
					log.Println("error getting URL:", u)
					log.Println(err)
					continue
				}

				buff := &bytes.Buffer{}
				_, err = io.Copy(buff, res.Body)
				if err != nil {
					log.Println("error reading data from URL:", u)
					log.Println(err)
					res.Body.Close()
					continue
				}

				res.Body.Close()

				contentType := res.Header.Get("Content-Type")
				if contentType == "" {
					contentType = "?"
				}

				log.Printf("Got %d bytes (%s) from %s", buff.Len(), contentType, u)

				fetchedQueue <- &fetchedURL{url: u, data: buff, contentType: contentType}
			}
		}()
	}

	return fetchedQueue
}

func (f fetchedURL) isTextHTML() bool {
	return strings.Contains(strings.ToLower(f.contentType), "text/html")
}

func (f fetchedURL) isJpegImage() bool {
	return strings.Contains(strings.ToLower(f.contentType), "image/jpeg")
}

func (f fetchedURL) isPngImage() bool {
	return strings.Contains(strings.ToLower(f.contentType), "image/png")
}
