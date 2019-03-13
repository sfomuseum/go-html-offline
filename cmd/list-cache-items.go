package main

import (
	"flag"
	"github.com/sfomuseum/go-html-offline"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/walk"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {

	mode := flag.String("mode", "file", "Indicate how command line arguments should be interpreted. Valid options are: files, directory.")
	validate := flag.Bool("validate", false, "...")

	var urls flags.MultiString
	flag.Var(&urls, "url", "One or more URLs to append to the service worker cache list")

	flag.Parse()

	opts := offline.DefaultServiceWorkerOptions()
	opts.CacheURLs = urls

	items := new(sync.Map)

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

			cache, err := offline.CacheListFromFile(path, opts)

			if err != nil {
				return err
			}

			items.Store(path, cache)
			return nil
		}

		for _, path := range flag.Args() {

			err := walk.Walk(path, cb)

			if err != nil {
				log.Fatal(err)
			}
		}

	case "file":

		for _, path := range flag.Args() {

			cache, err := offline.CacheListFromFile(path, opts)

			if err != nil {
				log.Fatal(err)
			}

			items.Store(path, cache)
		}

	case "url":

		for _, url := range flag.Args() {

			cache, err := offline.CacheListFromURL(url, opts)

			if err != nil {
				log.Fatal(err)
			}

			items.Store(url, cache)
		}

	default:
		log.Fatal("Invalid -mode")
	}

	items.Range(func(key interface{}, value interface{}) bool {

		uri := key.(string)
		cache := value.([]string)

		log.Println(uri)

		for _, u := range cache {
			log.Println(u)
		}

		return true
	})

	if *validate {

		to_validate := new(sync.Map)

		items.Range(func(key interface{}, value interface{}) bool {

			cache := value.([]string)

			for _, u := range cache {

				if strings.HasPrefix(u, "http") {
					to_validate.Store(u, true)
				}
			}

			return true
		})

		to_validate.Range(func(key interface{}, value interface{}) bool {

			uri := key.(string)

			rsp, err := http.Head(uri)

			if err != nil {
				log.Println("ERROR", uri, err)
			} else if rsp.StatusCode != 200 {
				log.Println("ERROR", uri, rsp.Status)
			} else {
				log.Println("OK", uri)
			}

			return true
		})

	}
}
