package main

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-html-offline"	
	"github.com/sfomuseum/go-html-offline/http"
	"log"
	gohttp "net/http"
)

func main() {

	cache_name := flag.String("cache-name", "network-or-cache", "The name for your browser/service worker cache.")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	flag.Parse()

	opts := offline.DefaultServiceWorkerOptions()
	opts.CacheName = *cache_name
	
	mux := gohttp.NewServeMux()

	ping_handler, err := http.PingHandler()

	if err != nil {
		log.Fatal(err)
	}

	inventory_handler, err := http.InventoryHandler(opts)

	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/ping", ping_handler)
	mux.Handle("/", inventory_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening for requests on %s\n", endpoint)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}
}
