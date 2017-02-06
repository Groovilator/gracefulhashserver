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

type GracefulListener struct {
	net.Listener
	wg *sync.WaitGroup
}

type GracefulConn struct {
	net.Conn
	wg *sync.WaitGroup
}

func NewGracefulListener(l net.Listener) (*GracefulListener) {
	return &GracefulListener{Listener: l, wg: &sync.WaitGroup{}}
}

func NewGracefulConn(wg *sync.WaitGroup, conn net.Conn) *GracefulConn {
	return &GracefulConn{Conn: conn, wg: wg}
}

func (gl *GracefulListener) Accept() (net.Conn, error) {
	conn, err := gl.Listener.Accept()
	if err == nil {
		sl.wg.Add(1)
	}
	return NewGracefulConn(gl.wg, conn), err
}

func (gconn *GracefulConn) Close() error {
	gconn.wg.Done()
	return gconn.Conn.Close()
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received hash request")

	r.ParseForm()
	pass := r.Form.Get("password")

	time.Sleep(5 * time.Second)

	sha_512 := sha512.New()
	sha_512.Write([]byte(pass))

	fmt.Fprintf(w, base64.StdEncoding.EncodeToString(sha_512.Sum(nil)))
	fmt.Println("Returning hash")
}

func main() {
	clumsyListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	gl := NewGracefulListener(clumsyListener)

	http.HandleFunc("/", hashHandler)
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("*Shutdown Started*")
		fmt.Printf("Closing listener\n")
		gl.Close()
	})
	server := http.Server{}

	fmt.Println("*Serving HTTP*")
	server.Serve(gl)

	fmt.Println("Waiting on pending responses")
	gl.wg.Wait()
}
