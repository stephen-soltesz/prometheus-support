package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type MultiString []string

func (m *MultiString) String() string {
	return strings.Join(*m, ",")
}

func (m *MultiString) Set(value string) error {
	*m = append(*m, value)
	return nil
}

// auth: none, env, user/pass.
// source data: url, project id
// dest: output file name.

func main() {
	client := http.Client{
		Timeout: time.Minute,
	}
	src := &MultiString{}
	dest := &MultiString{}
	var refresh time.Duration

	flag.Var(src, "source", "URL of a source.")
	flag.Var(dest, "dest", "Path to output.")
	flag.DurationVar(&refresh, "refresh", time.Minute, "Number of seconds between refreshing.")
	flag.Parse()

	if len(*src) != len(*dest) {
		log.Fatalf("Different lengths: src:%d dest:%d", len(*src), len(*dest))
	}

	var start time.Time

	// Only sleep as long as we need to, before starting a new iteration.
	for ; ; time.Sleep(refresh - time.Since(start)) {
		start = time.Now()
		log.Printf("Starting a new round at: %s", start)

		for i, source := range *src {
			log.Printf("%s -> %s\n", source, (*dest)[i])
			resp, err := client.Get(source)
			if err != nil {
				log.Println("Error:", err)
				continue
			}
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			err = ioutil.WriteFile((*dest)[i], data, 0644)
			if err != nil {
				log.Println("Error:", err)
				continue
			}
		}
	}
}
