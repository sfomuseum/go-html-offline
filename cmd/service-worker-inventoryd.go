package main

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-html-offline"
	"github.com/sfomuseum/go-html-offline/http"
	"github.com/sfomuseum/go-html-offline/server"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"	
	"log"
	gohttp "net/http"
	gourl "net/url"
	"strings"
)

func main() {

	cache_name := flag.String("cache-name", "network-or-cache", "The name for your browser/service worker cache.")
	var scheme = flag.String("scheme", "http", "The protocol scheme to use for the server. Valid options are: http, lambda.")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on.")
	var port = flag.Int("port", 8080, "The port number to listen for requests on.")
	var root = flag.String("root", "", "A valid URL to fetch subrequests from.")
	var path = flag.String("path", "/", "The path (URL) for handling requests.")
	var cors = flag.String("cors", "", "Set the following CORS access-control header.")
	var logging = flag.Bool("logging", false, "Log requests (to STDOUT).")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("INVENTORYD")

	if err != nil {
		log.Fatal(err)
	}
	
	if *root == "" {
		log.Fatal("Missing root")
	}

	_, err = gourl.Parse(*root)

	if err != nil {
		log.Fatal(err)
	}

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

	endpoint := fmt.Sprintf("%s://%s:%d", *scheme, *host, *port)
	url, err := gourl.Parse(endpoint)

	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer(url)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(mux)

	if err != nil {
		log.Fatal(err)
	}
}
