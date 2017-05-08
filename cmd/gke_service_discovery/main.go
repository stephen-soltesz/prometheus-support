// Note: this works with the code within the 201705 client-go repo.
package main

import (
	"flag"
	"log"
	"time"
)

var (
	project     = flag.String("project", "", "GCP project name.")
	gkeTargets  = flag.String("gke-targets", "targets.json", "Write targets configuration to given filename.")
	aefTargets  = flag.String("aef-targets", "targets.json", "Write targets configuration to given filename.")
	httpTargets = flag.String("http-targets", "targets.json", "Write targets configuration to given filename.")
	// scopeList   = flag.String("scopes", strings.Join(defaultScopes, ","), "Comma separated list of scopes to use.")
	refresh = flag.Duration("refresh", time.Minute, "Number of seconds between refreshing.")
)

type Target interface {
	Collect() error
	Save(output string) error
}

func main() {
	flag.Parse()
	start := time.Now()

	gkeSource, err := NewGKESource("mlab-sandbox")
	if err != nil {
		panic(err)
	}
	err = gkeSource.Collect()
	if err != nil {
		panic(err)
	}
	err = gkeSource.Save("output.json")
	if err != nil {
		panic(err)
	}
	log.Println(time.Since(start))
	if err != nil {
		panic(err)
	}
}
