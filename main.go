package main

import (
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
)

var cluster Cluster

func LoadConfig(){
    cluster.ID = os.Getenv("CLUSTER_ID")
}

func main(){
    flag.Parse()

    LoadConfig()
    log.Printf( "Returning creating projects under the %s cluster", cluster.ID )

    http.HandleFunc("/", Router)
    http.ListenAndServe(":8000", nil)
}
