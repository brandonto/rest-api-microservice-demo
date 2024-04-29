package main

import (
    //"fmt"
    "net/http"
)

func main() {
    r := NewRouter()

    http.ListenAndServe(":12345", r)
}
