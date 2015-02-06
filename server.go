package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
)

var cache map[[16]byte]*http.Response = make(map[[16]byte]*http.Response)

func handler(w http.ResponseWriter, r *http.Request) {
	var username string = r.URL.Query()["username"][0]
	var path string = r.URL.RequestURI()
	var cacheKey []byte = []byte(username + path)

	if resp, ok := cache[md5.Sum(cacheKey)]; ok {
		println("cache hit")
		fmt.Fprintf(w, resp.Status)
	} else {
		println("cache miss")
		res, err := http.Get("http://www.google.com/")
		if err != nil {
			log.Fatal(err)
		}
		cache[md5.Sum(cacheKey)] = res
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
