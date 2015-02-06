package main

import (
	"crypto/md5"
	"net/http"
)

var cache map[[16]byte]*http.Response

func handler(w http.ResponseWriter, r *http.Request) {
	var username string = r.URL.Query()["username"][0]
	var path string = r.URL.RequestURI()
	var cacheKey []byte = []byte(username + path)
	cache[md5.Sum(cacheKey)] = path
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
