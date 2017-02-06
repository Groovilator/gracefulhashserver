package main

import (
    "strings"
    "testing"
    "net/http"
    "io/ioutil"
    "fmt"
)

func TestKnownOutput(t *testing.T) {
    go main()

    resp, err := http.Get("http://localhost:8080/?password=angryMonkey")
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        t.Fatal(err)
    }

    if !strings.Contains(string(body), "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==") {
        fmt.Print(string(body))
        t.Errorf("Hash response doesn't match:\n%s", body)
    }

    _, _ = http.Get("http://localhost:8080/shutdown")
}
