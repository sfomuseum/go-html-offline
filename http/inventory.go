package http

import (
	"github.com/sfomuseum/go-html-offline"
	"io/ioutil"
	gohttp "net/http"
)

func InventoryHandler(opts *offline.ServiceWorkerOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		url := req.URL

		// test url here...

		rsp2, err := gohttp.Get(url.String())

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		defer rsp2.Body.Close()

		err = offline.AddServiceWorker(rsp2.Body, ioutil.Discard, rsp, opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
