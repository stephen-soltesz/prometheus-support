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
	var start time.Time

	// Only sleep as long as we need to, before starting a new iteration.
	for ; ; time.Sleep(*refresh - time.Since(start)) {
		start = time.Now()

		gkeSource, err := NewGKESource(*project)
		if err != nil {
			log.Printf("Failed to get authenticated client: %s", err)
			continue
		}
		err = gkeSource.Collect()
		if err != nil {
			log.Printf("Failed to get collect targets: %s", err)
			continue
		}
		err = gkeSource.Save(*gkeTargets)
		if err != nil {
			log.Printf("Failed to save to %s: %s", *gkeTargets, err)
			continue
		}
		log.Println(time.Since(start))
	}
}
