package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	recovered := true

	ticker := time.NewTicker(time.Second * time.Duration(cfg.Main.Interval))

	for now := range ticker.C {
		func() {
			resp, err := http.Get(cfg.Main.URL)
			defer resp.Body.Close()

			if err != nil && strings.Contains(err.Error(), "redirects") {
				if recovered {
					recovered = false
					writeit(cfg.Main.Filename, now.Format(time.RFC1123), strconv.Itoa(resp.StatusCode))
					fmt.Println(resp.StatusCode, "at", now, " - Updating", cfg.Main.Filename, "and waiting for", (cfg.Main.Wait), "seconds")
				} else {
					fmt.Println(resp.StatusCode, "at", now, " - Not recovered from previous failure, trying again in", (cfg.Main.Wait), "seconds")
				}
				time.Sleep(time.Duration(cfg.Main.Wait) * time.Second)
			} else if err != nil {
				fmt.Println("Unknown error while trying to receive", cfg.Main.URL, ". Waiting for", (cfg.Main.Wait), "seconds")
				time.Sleep(time.Duration(cfg.Main.Wait) * time.Second)
			} else {
				recovered = true
				fmt.Println(resp.StatusCode, "at", now)
			}
		}()
	}
}
