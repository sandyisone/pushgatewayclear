package main

import (
	"com.as/pgclear/work"
	"flag"
	"log"
	"time"
)

var (
	ttl = flag.Duration("ttl", 1*time.Hour, "TTL for clearing the expired metrics, the default is 1 hours.")
	url = flag.String("url", "", "url of Pushgateway")
)

func main() {

	log.Println("pushgatewayclear version 1.00 build21.03.20.11")

	flag.Parse()

	if *url == "" {
		log.Println("url of Pushgateway is null")
		return
	}


	work.Work(*url, *ttl)

}
