package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

func makeHashHandler(wg *sync.WaitGroup) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()
		fmt.Println("Received Hash Request")

		r.ParseForm()
		x := r.Form.Get("password")

		time.Sleep(5 * time.Second)

		sha_512 := sha512.New()
		sha_512.Write([]byte(x))

		fmt.Fprintf(w, base64.URLEncoding.EncodeToString(sha_512.Sum(nil)))
		fmt.Println("Returning Hash Response")
	}
}

func main() {
	clumsyListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	hashHandler := makeHashHandler(&wg)

	http.HandleFunc("/", hashHandler)
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("*Shutdown Started*")
		fmt.Printf("Closing listener\n")
		clumsyListener.Close()
	})
	server := http.Server{}

	fmt.Println("*Serving HTTP*")
	server.Serve(clumsyListener)

	fmt.Println("Waiting on pending responses")
	wg.Wait()
}
