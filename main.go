package main

import (
    "flag"
    "log"
    "net/http"
    "os"
)

var amqpUri *string
var cluster Cluster
var mode *string

func LoadConfig(){
    cluster.ID = os.Getenv("CLUSTER_ID")
}

func main(){
    amqpUri = flag.String("amqp-uri", "amqp://guest:guest@localhost:5671/vhost", "amqp uri when in rabbit mode")
    mode = flag.String("mode", "rest", "Mode in which to listen [rest|rabbit]")
    flag.Parse()

    LoadConfig()
    log.Printf( "Returning creating projects under the %s cluster", cluster.ID )

    switch *mode {
    case "rest":
        http.HandleFunc("/", Router)
        http.ListenAndServe(":8000", nil)
    case "rabbit":
        //
        _, err := NewConsumer(*amqpUri, "something", "somethingelse")
        if err != nil {
            log.Fatal(err.Error())
        }
    default:
        log.Fatalf("Mode %s is invalid", mode)
    }
}
