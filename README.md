# gracefulhashserver

Simple http server that accepts connections over port 8080 and provides a 
sha512 hash if a "password" parameter is provided. Additionally, a graceful
shutdown endpoint is exposed that stops the http server, refuses any 
subsequent requests, and exits.

## Endpoints:
**"/"** - Main endpoint that provides access to the sha512 hashing function. If 
      a "password" parameter is provided a base64 encoded sha512 hash of the
      string is returned. If the "password" parameter is not provided, the
      encoded hash of an empty string ("") will be returned.
      
   Ex. http://localhost:8080/?password=angryMonkey
      
**"/shutdown"** - Graceful shutdown endpoint. When the server receives the request
              it closes the http listener, refuses any new connections, and 
              waits for any requests that were sent before the shutdown to
              respond before ending execution.
   
   Ex. http://localhost:8080/shutdown
              
## Build/Run
While in the source directory with go installed on the executing machine, run
`go build main.go`
followed by
`./main`
to build and run. Alternatively,
`go run main.go`
will create a new build and run it.

For tests, run
`go test`
in the source directory.
