package main

// PLEASE RENAME ME...

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-html-offline"
	_ "github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/walk"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

func main() {

	mode := flag.String("mode", "file", "Indicate how command line arguments should be interpreted. Valid options are: files, directory.")

	flag.Parse()

	opts := offline.DefaultServiceWorkerOptions()

	urls_map := new(sync.Map)

	catalog := func(path string) error {

		cache, err := offline.CacheListFromFile(path, opts)

		if err != nil {
			return err
		}

		for _, url := range cache {
			urls_map.Store(url, true)
		}

		return nil
	}

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

			return catalog(path)
		}

		for _, path := range flag.Args() {

			err := walk.Walk(path, cb)

			if err != nil {
				log.Fatal(err)
			}
		}

	case "file":

		for _, path := range flag.Args() {

			err := catalog(path)

			if err != nil {
				log.Fatal(err)
			}
		}

	default:
		log.Fatal("Invalid -mode")
	}

	urls := make([]string, 0)

	urls_map.Range(func(key interface{}, value interface{}) bool {

		url := key.(string)
		urls = append(urls, url)
		return true
	})

	sort.Strings(urls)

	for _, u := range urls {
		fmt.Println(u)
	}
}
