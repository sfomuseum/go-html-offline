package offline

import ()

var sw string
var sw_init string

func init() {

	sw = `
// this file was generated by robots on {{ .Date }}
// https://github.com/sfomuseum/go-html-offline

var CACHE = '{{ .CacheName }}';

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
	{{ range  $uri := .ToCache }}'{{ $uri }}',
	{{ end }}
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
}`

	sw_init = `
// this code was added by robots on {{ .Date }}
// https://github.com/sfomuseum/go-html-offline

window.addEventListener("load", function load(event){
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('{{ .ServiceWorkerURL }}').then(function(registration) {
    console.log('Service worker registration succeeded:', registration);
  }, /*catch*/ function(error) {
    console.log('Service worker registration failed:', error);
  });
} else {
  console.log('Service workers are not supported.');
}
}, false);`

}
