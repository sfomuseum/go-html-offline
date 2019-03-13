package http

import (
	"bufio"
	"bytes"
	"github.com/sfomuseum/go-html-offline"
	"io/ioutil"
	"log"
	gohttp "net/http"
	gourl "net/url"
	"strconv"
	"strings"
)

type InventoryOptions struct {
	Root    string
	CORS    string
	Path    string
	Logging bool
}

func InventoryHandler(inv_opts *InventoryOptions, sw_opts *offline.ServiceWorkerOptions) (gohttp.Handler, error) {

	root, err := gourl.Parse(inv_opts.Root)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		url := req.URL

		url.Scheme = root.Scheme
		url.Host = root.Host

		path := strings.Replace(url.Path, inv_opts.Path, "", 1)
		url.Path = path

		if inv_opts.Logging {
			log.Printf("Fetch '%s'\n", url)
		}

		rsp2, err := gohttp.Get(url.String())

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		defer rsp2.Body.Close()

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		err = offline.AddServiceWorker(rsp2.Body, ioutil.Discard, wr, sw_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		wr.Flush()

		data := buf.Bytes()
		clen := len(data)

		rsp.Header().Set("Content-Length", strconv.Itoa(clen))
		rsp.Header().Set("Content-Type", "text/javascript")

		if inv_opts.CORS != "" {
			rsp.Header().Set("Access-Control-Allow-Origin", inv_opts.CORS)
		}

		rsp.Write(data)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
