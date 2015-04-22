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

func main() {
	var cfg config
	gcfg.ReadFileInto(&cfg, "config.gcfg")

	// url := "http://localhost:9090"
	// interval := time.Duration(5000)
	// naptime := time.Duration(30000)
	// filename := "web.config"

	ticker := time.NewTicker(time.Millisecond * time.Duration(cfg.Main.Interval))
	for now := range ticker.C {
		resp, err := http.Get(cfg.Main.URL)
		//	var status string

		if (err != nil) || (resp.StatusCode > 299 && resp.StatusCode < 400) {
			f, err := os.OpenFile(cfg.Main.Filename, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			if _, err = f.WriteString("<!-- redirect loop detected, restarted at: " + now.Format(time.RFC1123) + " Status: " + strconv.Itoa(resp.StatusCode) + " -->\n"); err != nil {
				panic(err)
			}
			fmt.Println(resp.StatusCode, "at", now, " - Updating", cfg.Main.Filename, "and sleeping for", (cfg.Main.Wait / 1000), "seconds")

			time.Sleep(time.Duration(cfg.Main.Wait) * time.Millisecond)

		} else {
			fmt.Println(resp.StatusCode, "at", now)
		}
	}
}
