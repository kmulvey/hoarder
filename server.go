package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CachedReq struct {
	header http.Header
	body   io.ReadCloser
}

var cache map[[16]byte]CachedReq = make(map[[16]byte]CachedReq)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/dump", dumpHandler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// not checking errors on query .. need something like username, ok := Query()["key"]
	var username string = r.URL.Query()["key"][0]
	var path string = r.URL.RequestURI()
	var cacheKey []byte = []byte(username + path)

	if resp, ok := cache[md5.Sum(cacheKey)]; ok {
		println("cache hit")
		// send it back
		copyHeaders(w.Header(), r.Header)
		io.Copy(w, resp.body)
	} else {
		println("cache miss")
		res, err := http.Get("http://localhost")
		if err != nil {
			log.Fatal(err)
		}
		// serve it back
		copyHeaders(w.Header(), res.Header)
		io.Copy(w, res.Body)

		// populate the cache
		cache[md5.Sum(cacheKey)] = CachedReq{
			res.Header,
			res.Body,
		}
	}
}

// dumpHandler is supposed to dump the cache but its broken
// at the moment
func dumpHandler(w http.ResponseWriter, r *http.Request) {
	var keys string
	for k, _ := range cache {
		keys += string(k[:16]) + " "
	}
	fmt.Fprintf(w, keys)
}

// copyHeaders overwrites headers and it sucks becasue
// we are looping on evey req
func copyHeaders(dst, src http.Header) {
	for k, w := range src {
		for _, v := range w {
			dst.Add(k, v)
		}
	}
}
