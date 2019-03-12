package offline

// https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API/Using_Service_Workers
// https://hacks.mozilla.org/2015/11/offline-service-workers/
// https://serviceworke.rs/strategy-network-or-cache_service-worker_doc.html
// https://developer.mozilla.org/en-US/docs/Web/API/Cache

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	_ "log"
	"strings"
	"text/template"
)

type ServiceWorkerOptions struct {
	CacheName        string
	ServiceWorkerURL string
}

func DefaultServiceWorkerOptions() *ServiceWorkerOptions {

	opts := ServiceWorkerOptions{
		CacheName:        "network-or-cache",
		ServiceWorkerURL: "sw.js",
	}

	return &opts
}

func AddServiceWorker(in io.Reader, html_wr io.Writer, serviceworker_wr io.Writer, opts *ServiceWorkerOptions) error {

	sw_t, err := template.New("service-worker").Parse(sw)

	if err != nil {
		return err
	}

	init_t, err := template.New("service-worker-init").Parse(sw_init)

	if err != nil {
		return err
	}

	doc, err := html.Parse(in)

	if err != nil {
		return err
	}

	type ServiceWorkerVars struct {
		CacheName string
		ToCache   []string
	}

	type ServiceWorkerInitVars struct {
		ServiceWorkerURL string
	}

	to_cache := []string{
		"",
		"index.html",
	}

	var callback func(node *html.Node, writer io.Writer)

	callback = func(n *html.Node, w io.Writer) {

		if n.Type == html.ElementNode {

			switch n.Data {

			case "head":

				for c := n.FirstChild; c != nil; c = c.NextSibling {

					if c.Type != html.ElementNode || c.Data != "script" {
						continue
					}

					script := attrs2map(c.Attr...)

					_, ok := script["x-service-worker"]

					if ok {
						n.RemoveChild(c)
					}
				}

				vars := ServiceWorkerInitVars{
					ServiceWorkerURL: opts.ServiceWorkerURL,
				}

				var buf bytes.Buffer
				wr := bufio.NewWriter(&buf)

				err := init_t.Execute(wr, vars)

				if err != nil {
					// log.Println(err)
					return
				}

				wr.Flush()

				script_type := html.Attribute{"", "type", "text/javascript"}
				script_rel := html.Attribute{"", "x-service-worker", "true"}

				script := html.Node{
					Type:      html.ElementNode,
					DataAtom:  atom.Script,
					Data:      "script",
					Namespace: "",
					Attr:      []html.Attribute{script_type, script_rel},
				}

				body := html.Node{
					Type: html.TextNode,
					Data: string(buf.Bytes()),
				}

				script.AppendChild(&body)
				n.AppendChild(&script)

			case "img":

				for _, attr := range n.Attr {

					if attr.Key == "src" {
						to_cache = append(to_cache, attr.Val)
						break
					}
				}

			case "link":

				link := attrs2map(n.Attr...)

				rel, rel_ok := link["rel"]
				href, href_ok := link["href"]

				if rel_ok && href_ok && rel == "stylesheet" {
					to_cache = append(to_cache, href)
				}

			case "script":

				script := attrs2map(n.Attr...)

				script_type, script_type_ok := script["type"]
				src, src_ok := script["src"]

				if script_type_ok && src_ok && script_type == "text/javascript" {
					to_cache = append(to_cache, src)
				}

			case "source":

				// <picture> uses <source srcset="...">
				// <video> uses <source src="...">

				source := attrs2map(n.Attr...)

				srcset, srcset_ok := source["srcset"]

				if srcset_ok {
					to_cache = append(to_cache, srcset)
				}

				src, src_ok := source["src"]

				if src_ok {
					to_cache = append(to_cache, src)
				}

			default:
				// pass
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			callback(c, html_wr)
		}
	}

	callback(doc, html_wr)

	for idx, uri := range to_cache {

		if !strings.HasPrefix(uri, "/") {
			to_cache[idx] = fmt.Sprintf("./%s", uri)
		}
	}

	vars := ServiceWorkerVars{
		CacheName: opts.CacheName,
		ToCache:   to_cache,
	}

	err = sw_t.Execute(serviceworker_wr, vars)

	if err != nil {
		return err
	}

	return html.Render(html_wr, doc)
}

func attrs2map(attrs ...html.Attribute) map[string]string {

	attrs_map := make(map[string]string)

	for _, a := range attrs {
		attrs_map[a.Key] = a.Val
	}

	return attrs_map
}
