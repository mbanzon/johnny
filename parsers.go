package main

import (
	"bytes"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func setupParsers(count int, inputQueue chan []string, parseQueue chan *fetchedURL) {
	log.Println("Setting up parsers:", count)

	for i := 0; i < count; i++ {
		go parser(inputQueue, parseQueue)
	}
}

func parser(inputQueue chan []string, parseQueue chan *fetchedURL) {
	for {
		f := <-parseQueue

		log.Printf("Parsing %d bytes as %s", f.data.Len(), f.contentType)

		if f.isTextHTML() {
			links, err := parseHTMLForLinks(f.url, f.data)
			if err != nil {
				log.Println("error parsing HTML from URL:", f.url)
				log.Println(err)
			}

			inputQueue <- links
		} else if f.isJpegImage() {
			// TODO: implement
		} else if f.isPngImage() {
			// TODO: implement
		} else {
			log.Println("Unknown content type:", f.contentType)
		}
	}
}

func parseHTMLForLinks(base string, data *bytes.Buffer) ([]string, error) {
	doc, err := html.Parse(data)
	if err != nil {
		return nil, err
	}

	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			link := ""
			if n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						link = createLink(base, a.Val)
						break
					}
				}
			} else if n.Data == "img" {
				for _, a := range n.Attr {
					if a.Key == "src" {
						link = createLink(base, a.Val)
						break
					}
				}
			}

			if link != "" {
				links = append(links, link)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return links, nil
}

func createLink(base, link string) string {
	if strings.HasPrefix(link, "mailto") {
		return ""
	}

	if strings.HasPrefix(link, "ftp") {
		return ""
	}

	if strings.HasPrefix(link, "http") {
		return link
	}

	if strings.HasPrefix(link, "data:") {
		return "" // TODO: handle image data
	}

	if strings.HasPrefix(link, "//") {
		return "https:" + link
	}

	if strings.HasPrefix(link, "/") {
		slashCount := strings.Count(base, "/")
		if slashCount == 2 {
			return base + link
		} else if slashCount == 3 {
			if strings.HasSuffix(base, "/") {
				return base[:len(base)-1] + link
			}
			return base + link
		} else {
			index := 0
			for i := 0; i < 3; i++ {
				sub := base[index:]
				index += strings.Index(sub, "/") + 1
			}
			return base[:index-1] + link
		}
	}
	return base[:strings.LastIndex(base, "/")+1] + link
}
