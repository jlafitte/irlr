package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"code.google.com/p/gcfg"
)

type config struct {
	Main struct {
		URL      string
		Interval int
		Wait     int
		Filename string
	}
}

func writeit(file string, when string, what string) {
	f, _ := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	f.WriteString("<!-- redirect loop detected, restarted at: " + when + " Status: " + what + " -->\n")
}

func main() {
	var cfg config
	gcfg.ReadFileInto(&cfg, "config.gcfg")

	ticker := time.NewTicker(time.Second * time.Duration(cfg.Main.Interval))
	for now := range ticker.C {
		resp, err := http.Get(cfg.Main.URL)
		//	var status string

		if (err != nil) || (resp.StatusCode > 299 && resp.StatusCode < 400) {
			writeit(cfg.Main.Filename, now.Format(time.RFC1123), strconv.Itoa(resp.StatusCode))
			fmt.Println(resp.StatusCode, "at", now, " - Updating", cfg.Main.Filename, "and sleeping for", cfg.Main.Wait, "seconds")
			time.Sleep(time.Duration(cfg.Main.Wait) * time.Second)
		} else {
			fmt.Println(resp.StatusCode, "at", now)
		}
	}
}
