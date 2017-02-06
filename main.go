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

// Listener wrapper to hold server state information
type GracefulListener struct {
	net.Listener
	wg *sync.WaitGroup
}

// Conn wrapper to hold server state information
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

// Listener.Accept() implementation that add's to the server's waitgroup if a 
// connection is successful created
func (gl *GracefulListener) Accept() (net.Conn, error) {
	conn, err := gl.Listener.Accept()
	if err == nil {
		sl.wg.Add(1)
	}
	return NewGracefulConn(gl.wg, conn), err
}

// Conn.Close() implemenation that decrements the server's waitgroup when a 
// connection is closed d
func (gconn *GracefulConn) Close() error {
	gconn.wg.Done()
	return gconn.Conn.Close()
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received hash request")

	// parse "password" parameter, pass = "" if not specified
	r.ParseForm()
	pass := r.Form.Get("password")

	// block for 5 seconds
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
		// Close listener (in turn stopping the server)
		gl.Close()
	})
	server := http.Server{}

	fmt.Println("*Serving HTTP*")
	server.Serve(gl)
            
	fmt.Println("Waiting on pending responses")
	gl.wg.Wait()
}

