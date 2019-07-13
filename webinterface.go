package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // include pprof to do some memory profiling
	"net/url"
)

func startupServer(noServer, publicServer bool, queue chan []string) {
	if !noServer {
		addr := "localhost:5000"
		if publicServer {
			addr = ":5000"
		}

		log.Println("Setting up server...")

		setupHandlers(queue)

		log.Println("Starting server...")

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Println("server listen error:", err)
		}
	} else {
		log.Println("Not starting server...")
	}
}

func setupHandlers(queue chan []string) {
	http.HandleFunc("/queue", makeQueueHandler(queue))
	http.HandleFunc("/", rootHandler)
}

func makeQueueHandler(queue chan []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		var err error

		if err = r.ParseForm(); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		urlStr := r.FormValue("url")
		if urlStr == "" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		var u *url.URL

		if u, err = url.Parse(urlStr); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		go func() {
			queue <- []string{u.String()}
		}()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!doctype html>
<html>
	<head>
		<title>Spider - Web Scraper</title>
	</head>
	<body>
		<form method="post" action="/queue">
			<label>Queue URL:</label><br />
			<input type="url" name="url" required /><br />
			<input type="submit" value="Put in queue" />
		</form>
	</body>
</html>
`)
}
