package main

import (
	"github.com/diamondo25/msps-tracker/lib"
	"github.com/pkg/errors"
	"log"
	"time"
)

func main() {
	cfg := lib.NewConfig("config.yml")

	servers := cfg.Servers

	refreshInterval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		panic(errors.Wrap(err, "unable to parse interval"))
	}

	for _, ps := range servers {
		go func(ps lib.ServerConfig) {
			t := time.Tick(refreshInterval)
			for {
				start := time.Now()
				amount, err := ps.FetchAmount()
				if err != nil {
					log.Println("[", ps.Name, "] unable to fetch: ", err)
				} else {
					log.Println("[", ps.Name, "] Online count: ", amount)

					doc := &lib.ElasticSearchDocument{}
					doc.Timestamp = time.Now()
					doc.OnlineCount = amount
					doc.ServerName = ps.Name
					doc.ResponseTimeMilliseconds = int(time.Since(start) / time.Millisecond)

					if err := cfg.ElasticSearch.Write(doc); err != nil {
						log.Println("Unable to write to ES", err)
					}
				}

				<-t
			}
		}(ps)
	}

	select {}
}
