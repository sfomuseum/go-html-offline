package main

import (
	"flag"
	"github.com/sfomuseum/go-html-offline"
	"github.com/whosonfirst/walk"
	"log"
	"os"
	"strings"
)

func main() {

	cache_name := flag.String("cache-name", "network-or-cache", "The name for your browser/service worker cache.")
	sw_url := flag.String("server-worker-url", "sw.js", "The URI of the JavaScript service worker.")
	mode := flag.String("mode", "file", "Indicate how command line arguments should be interpreted. Valid options are: files, directory.")

	flag.Parse()

	opts := offline.DefaultServiceWorkerOptions()
	opts.CacheName = *cache_name
	opts.ServiceWorkerURL = *sw_url

	switch *mode {

	case "directory":

		cb := func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".html") {
				return nil
			}

			return offline.AddServiceWorkerToFile(path, opts)
		}

		for _, path := range flag.Args() {

			err := walk.Walk(path, cb)

			if err != nil {
				log.Fatal(err)
			}
		}

	case "file":

		for _, path := range flag.Args() {

			err := offline.AddServiceWorkerToFile(path, opts)

			if err != nil {
				log.Fatal(err)
			}
		}

	default:
		log.Fatal("Invalid -mode")
	}

}
