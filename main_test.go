package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.google.com", 301)
}

func main() {
	url := "http://localhost:9090"
	interval := time.Duration(5000)
	naptime := time.Duration(30000)
	filename := "web.config"

	go func(msg string) {
		http.HandleFunc("/", redirect)
		http.ListenAndServe(":9090", nil)
	}("Starting web server")

	ticker := time.NewTicker(time.Millisecond * interval)
	for now := range ticker.C {
		resp, err := http.Get(url)
		//	var status string

		if (err != nil) || (resp.StatusCode > 299 && resp.StatusCode < 400) {
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			if _, err = f.WriteString("<!-- redirect loop detected, restarted at: " + now.Format(time.RFC1123) + " Status: " + strconv.Itoa(resp.StatusCode) + " -->\n"); err != nil {
				panic(err)
			}
			fmt.Println(resp.StatusCode, "at", now, " - Attempting to restart app pool and sleeping for", naptime, "seconds")

			time.Sleep(naptime * time.Millisecond)

		} else {
			fmt.Println(resp.StatusCode, "at", now)
		}
	}
}
