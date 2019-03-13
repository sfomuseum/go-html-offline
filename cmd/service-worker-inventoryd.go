package main

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-html-offline"
	"github.com/sfomuseum/go-html-offline/http"
	"log"
	gohttp "net/http"
	"strings"
)

func main() {

	cache_name := flag.String("cache-name", "network-or-cache", "The name for your browser/service worker cache.")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on.")
	var port = flag.Int("port", 8080, "The port number to listen for requests on.")
	var root = flag.String("root", "", "A valid URL to fetch subrequests from.")
	var path = flag.String("path", "/", "The path (URL) for handling requests.")
	var cors = flag.String("cors", "", "Set the following CORS access-control header.")
	var logging = flag.Bool("logging", false, "Log requests (to STDOUT).")

	flag.Parse()

	if !strings.HasSuffix(*path, "/") {
		*path = fmt.Sprintf("%s/", *path)
	}

	sw_opts := offline.DefaultServiceWorkerOptions()
	sw_opts.CacheName = *cache_name

	inv_opts := http.InventoryOptions{
		Root:    *root,
		Path:    *path,
		CORS:    *cors,
		Logging: *logging,
	}

	mux := gohttp.NewServeMux()

	ping_handler, err := http.PingHandler()

	if err != nil {
		log.Fatal(err)
	}

	inventory_handler, err := http.InventoryHandler(&inv_opts, sw_opts)

	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/ping", ping_handler)
	mux.Handle(*path, inventory_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening for requests on %s%s\n", endpoint, *path)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}
}
