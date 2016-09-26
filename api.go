package main

import (
    "net/http"
    "log"
    "io"
    "bytes"
    "fmt"
    "encoding/json"
)

var responseBody string
var status int
var errors bool
var resp Response

type Response struct {
    Body string
    Status int
}

func LogRequest(r *http.Request) {
    log.Printf( "%s :: %s %s",
        r.RemoteAddr,
        r.Method,
        r.URL.Path)
}

func Headers(w http.ResponseWriter) (http.ResponseWriter){
    w.Header().Set("Access-Control-Allow-Headers", "requested-with, Content-Type, origin, authorization, accept, client-security-token, cache-control, x-api-key")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Max-Age", "10000")
    w.Header().Set("Cache-Control", "no-cache")

    w.Header().Set("Content-Type", "application/json")

    return w
}

func Router(w http.ResponseWriter, r *http.Request){
    LogRequest(r)
    w = Headers(w)

    switch  {
    case r.URL.Path ==  "/" && r.Method == "POST":
        createProject(w,r)
    default:
        resp.Status = http.StatusNotFound
        resp.Body = "Not found"
        resp.respond(w)
    }
}

func createProject(w http.ResponseWriter, r *http.Request) {
    resp.Body, errors = NewProject(rcToString(r.Body))

    if errors == true {
        resp.Status = http.StatusInternalServerError
    } else {
        resp.Status = http.StatusOK
    }

    resp.respond(w)
}

func (r Response) respond (w http.ResponseWriter) {
    w.WriteHeader(r.Status)
    j,_ := json.Marshal(r)
    fmt.Fprintf(w, string(j))
}

func rcToString(rc io.ReadCloser) (s string){
    buf := new(bytes.Buffer)
    buf.ReadFrom(rc)
    s = buf.String()

    return
}
