package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CachedReq struct {
	header http.Header
	body   []byte
}

var cache map[[16]byte]CachedReq = make(map[[16]byte]CachedReq)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/dump", dumpHandler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var username string = r.URL.Query()["username"][0]
	var path string = r.URL.RequestURI()
	var cacheKey []byte = []byte(username + path)

	if resp, ok := cache[md5.Sum(cacheKey)]; ok {
		println("cache hit")
		w.Write(resp.body)
	} else {
		println("cache miss")
		res, err := http.Get("http://www.google.com/")
		if err != nil {
			log.Fatal(err)
		}
		body, bErr := ioutil.ReadAll(res.Body)
		if bErr != nil {
			log.Fatal(bErr)
		}
		cache[md5.Sum(cacheKey)] = CachedReq{
			res.Header,
			body,
		}
	}
}

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	var keys string
	for k, _ := range cache {
		keys += string(k[:16]) + " "
	}
	fmt.Fprintf(w, keys)
}
