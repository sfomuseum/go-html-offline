# go-html-offline

Tools for making HTML files service-worker (offline) ready.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

```
import (
	"github.com/sfomuseum/go-html-offline"
)

html_path := "/path/to/index.html"

opts := offline.DefaultServiceWorkerOptions()
offline.AddServiceWorkerToFile(html_path, opts)
```

The `AddServiceWorkerToFile` method is a helper method, working with atomic files, around the more abstract `AddServiceWorker` method that assumes `io.Reader` and `io.Writer` interfaces:

```
import (
	"github.com/sfomuseum/go-html-offline"
	"os"
)

html_path := "/path/to/index.html"

html_in, _ := os.Open(html_path)
html_out := os.Stdout			// update to point to service worker javascript
sw_out := os.Stdout			// the actual service worker javascript

opts := offline.DefaultServiceWorkerOptions()
offline.AddServiceWorker(html_in, html_out, sw_out, opts)
```

_Note that error handling has been removed for the sake of brevity._

The `AddServiceWorker` methods will update the HTML markup, reading from the `html_in` and writing to the `html_out` interfaces, to include the following JavaScript content:

```
<script type="text/javascript" x-service-worker="true">

// this code was added by robots on 2019-03-12T16:07:05-07:00
// https://github.com/sfomuseum/go-html-offline

window.addEventListener("load", function load(event){

if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('sw.js').then(function(registration) {
    console.log('Service worker registration succeeded:', registration);
  }, /*catch*/ function(error) {
    console.log('Service worker registration failed:', error);
  });
} else {
  console.log('Service workers are not supported.');
}

}, false);

</script>
```

It will also write the following JavaScript to the `sw_out` `io.Writer` interface. The array of URLs passed to the JavaScript `cache.addAll` method are determined from the body of the `html_in` HTML file.

```
// this file was generated by robots on 2019-03-12T16:07:05-07:00
// https://github.com/sfomuseum/go-html-offline

var CACHE = 'network-or-cache';

self.addEventListener('install', function(evt) {
  console.log('The service worker is being installed.');
  evt.waitUntil(precache());
});

self.addEventListener('fetch', function(evt) {
  console.log('The service worker is serving the asset.');
  evt.respondWith(fromNetwork(evt.request, 400).catch(function () {
    return fromCache(evt.request);
  }));
});

function precache() {
  return caches.open(CACHE).then(function (cache) {
    return cache.addAll([
	'./',
	'./index.html',
	...other assets in html_path
    ]);
  });
}

function fromNetwork(request, timeout) {
  return new Promise(function (fulfill, reject) {
    var timeoutId = setTimeout(reject, timeout);
    fetch(request).then(function (response) {
      clearTimeout(timeoutId);
      fulfill(response);
    }, reject);
  });
}

function fromCache(request) {
  return caches.open(CACHE).then(function (cache) {
    return cache.match(request).then(function (matching) {
      if (! matching){
	return Promise.reject('no-match');
      }
      return matching;
    });
  });
}
```

## URIs

URIs (to be passed to the `cache.addAll` JavaScript function) are derived from the following HTML elements:

* &lt;img src="{URI}" /&gt;
* &lt;link rel="stylesheet" href="{URI}" /&gt;
* &lt;script type="text/javascript" src="{URI}" /&gt;
* &lt;source srcset="{URI}" /&gt;
* &lt;source src="{URI}" /&gt;

## Tools

### add-service-worker

Add server worker JavaScript handlers for one or more HTML files.

```
./bin/add-service-worker -h
Usage of ./bin/add-service-worker:
  -cache-name string
    	The name for your browser/service worker cache. (default "network-or-cache")
  -mode string
    	Indicate how command line arguments should be interpreted. Valid options are: files, directory. (default "file")
  -server-worker-url string
    	The URI of the JavaScript service worker. (default "sw.js")
  -url value
    	One or more URLs to append to the service worker cache list
```

For example:

```
$> add-service-worker /path/to/index.html
$> ls -al /path/to/index.html
/path/to/index.html
/path/to/sw.js
```

If you just want to bulk process one or more folders full of `.html` files you would invoke `add-service-worker` with the `-mode directory` flag.

### service-worker-inventoryd

`service-worker-inventoryd` is an HTTP server that will fetch a URL and generate a "network-or-cache" style service worker JavaScript file for the assets listed in that page (URL) using the `offline.AddServiceWorker` method.

The server does not provide any access controls so if that's important to you then you will need to run this server behind something that does.

```
./bin/service-worker-inventoryd -h
Usage of ./bin/service-worker-inventoryd:
  -cache-name string
    	The name for your browser/service worker cache. (default "network-or-cache")
  -cors string
    	Set the following CORS access-control header.
  -host string
    	The hostname to listen for requests on. (default "localhost")
  -httptest.serve string
    	if non-empty, httptest.NewServer serves on this address and blocks
  -logging
    	Log requests (to STDOUT).
  -path string
    	The path (URL) for handling requests. (default "/")
  -port int
    	The port number to listen for requests on. (default 8080)
  -root string
    	A valid URL to fetch subrequests from.
  -scheme string
    	Valid options are: http, lambda. (default "http")
  -url value
    	One or more URLs to append to the service worker cache list	
```

For example:

```
$> ./bin/service-worker-inventoryd -prefix "/sw" -root "https://millsfield.sfomuseum.org"
2019/03/13 10:20:12 listening for requests on localhost:8080/sw/
```

And then in another terminal:

```
$> curl localhost:8080/sw/images/1377126959/

// this file was generated by robots on 2019-03-12T18:00:06-07:00
// https://github.com/sfomuseum/go-html-offline

... omitted for the sake of brevity

function precache() {
  return caches.open(CACHE).then(function (cache) {
    return cache.addAll([
	'./',
	'./index.html',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum-common/bootstrap.4.1.1.min.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum-common/leaflet.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum-common/sfomuseum.millsfield.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum-common/sfomuseum.millsfield.print.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum-common/sfomuseum.millsfield.placetypes.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum.collection.bootstrap.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum.collection.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum-common/sfomuseum.millsfield.leaflet.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum-common/sfomuseum.millsfield.media.css',
	'https://millsfield.sfomuseum.org/css/collection/sfomuseum.collection.media.css',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_c.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_z.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_n.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_k.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_b.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_c.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_z.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_n.jpg',
	'https://millsfield.sfomuseum.org/media/137/712/695/9/1377126959_fORVrjb1WIuTRkNHMDIO9TRDKqv0x3ce_c.jpg',
	'https://millsfield.sfomuseum.org/javascript/collection/sfomuseum.collection.maps.features.init.js',
	'https://millsfield.sfomuseum.org/javascript/collection/sfomuseum.collection.maps.results.init.js',
	
    ]);
  });
}

... omitted for the sake of brevity
```

## AWS (Lambda)

It is possible to run the `service-worker-inventoryd` daemon as Lambda function and accessed via an API Gateway endpoint.

The first step is to build the Lambda function to upload to AWS. This can be done using the handy `make lambda` Makefile target.

```
$> cd go-html-offline
$> make lambda
if test -d pkg; then rm -rf pkg; fi
if test -s src; then rm -rf src; fi
mkdir -p src/github.com/sfomuseum/go-html-offline
cp *.go src/github.com/sfomuseum/go-html-offline/
cp -r http src/github.com/sfomuseum/go-html-offline/
cp -r server src/github.com/sfomuseum/go-html-offline/
cp -r vendor/* src/
if test -f main; then rm -f main; fi
if test -f deployment.zip; then rm -f deployment.zip; fi
zip deployment.zip main
  adding: main (deflated 51%)
rm -f main
```

Now, upload `deployment.zip` to AWS. Your Lambda function (let's just say it's called `ServiceWorkerJS`) will need to be configured as follows:

* The runtime is `Go 1.x`
* The handler name is `main`
* Your functions executes as role with minimum `AWSLambdaExecute` permissions

You will need to set the following environment variables:

| Key | Value |
| --- | --- |
| INVENTORYD_SCHEME | lambda |
| INVENTORYD_ROOT | https://THE-HOST-YOU-WANT-GENERATE-SERVICE-WORKER-JS-FILES-FOR |

Any of the command line parameters can be passed to the Lambda function by setting environment variables as follows:

* Upper case the flag name and replace all spaces with the `_` character
* Prefix the environment variable with `INVENTORYD_`

To configure the API Gateway endpoint to invoke your Lambda function, head over to the API Gateway console and:

* Create a new resource
* In the `resource path` enter `{+proxy}` (name it whatever you want) - or just check the `Configure as proxy resource` button.
* Check the `Enable API Gateway CORS` button.
* Adjust the resource methods as necessary, in this case adding a `GET` method pointing to the `ServiceWorkerJS` lambda function.
* Remove the `ANY` and the `OPTIONS` methods.
* Deploy the API (with a new stage called `sw` or whatever you want).

Then you should be able to do something like this:

```
$> curl -s -v https://{ENDPOINT}.execute-api.{REGION}.amazonaws.com/sw/PATH-YOU-WANT-GENERATE-SERVICE-WORKER-JS-FILE-FOR > /dev/null

< HTTP/2 200 
< date: Wed, 13 Mar 2019 19:11:51 GMT
< content-type: text/javascript
< content-length: 3249
< x-amzn-requestid: db5105a7-45c3-11e9-9b01-01c7c06f5b41
< x-amzn-remapped-content-length: 3249
< x-amz-apigw-id: {GATEWAY_ID}
< x-amzn-trace-id: Root=1-5c8955f7-4b9b0280fecda0c0a0eeccc0;Sampled=0
```

## See also

* https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API/Using_Service_Workers
* https://hacks.mozilla.org/2015/11/offline-service-workers/
* https://serviceworke.rs/strategy-network-or-cache_service-worker_doc.html