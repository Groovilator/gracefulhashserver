package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func hashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received hash request")

	r.ParseForm()
	x := r.Form.Get("password")

	time.Sleep(5 * time.Second)

	sha_512 := sha512.New()
	sha_512.Write([]byte(x))

	fmt.Fprintf(w, base64.URLEncoding.EncodeToString(sha_512.Sum(nil)))
	fmt.Println("Returning hash")
}

func main() {
	http.HandleFunc("/", hashHandler)

	fmt.Println("*Serving HTTP*")
	http.ListenAndServe(":8080", nil)
}
