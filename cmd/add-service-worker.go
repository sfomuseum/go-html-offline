package main

// https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API/Using_Service_Workers
// https://hacks.mozilla.org/2015/11/offline-service-workers/
// https://serviceworke.rs/strategy-network-or-cache_service-worker_doc.html
// https://developer.mozilla.org/en-US/docs/Web/API/Cache

import (
	"flag"
	"github.com/facebookgo/atomicfile"
	"github.com/sfomuseum/go-html-offline"
	"log"
	"os"
	"path/filepath"
)

func main() {

	cache_name := flag.String("cache-name", "network-or-cache", "The name for your browser/service worker cache.")
	sw_url := flag.String("server-worker-url", "sw.js", "The URI of the JavaScript service worker.")

	flag.Parse()

	opts := offline.DefaultServiceWorkerOptions()
	opts.CacheName = *cache_name
	opts.ServiceWorkerURL = *sw_url

	for _, path := range flag.Args() {

		html_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		in, err := os.Open(html_path)

		if err != nil {
			log.Fatal(err)
		}

		root := filepath.Dir(html_path)
		sw_path := filepath.Join(root, *sw_url)

		html_out, err := atomicfile.New(html_path, 0644)

		if err != nil {
			log.Fatal(err)
		}

		sw_out, err := atomicfile.New(sw_path, 0644)

		if err != nil {
			html_out.Abort()
			log.Fatal(err)
		}

		err = offline.AddServiceWorker(in, html_out, sw_out, opts)

		if err != nil {
			html_out.Abort()
			sw_out.Abort()
			log.Fatal(err)
		}

		in.Close()

		html_out.Close()
		sw_out.Close()
	}

}
