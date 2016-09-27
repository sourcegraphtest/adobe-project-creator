package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "strings"
)

type Project struct {
    UUID string
    Name string
}
var project Project

func NewProject(body string)(statusString string, errors bool) {
    var statuses []string
    var status string
    var err error

    errors = false

    err = json.Unmarshal([]byte(body), &project)
    if err != nil {
        log.Println("error: ", err)
    }

    plproj,_ := ioutil.ReadFile("./empty.plproj")
    prproj,_ := ioutil.ReadFile("./empty.prproj")

    status, err = uploadData(cluster, project, "plproj", plproj)
    if err != nil {
        statuses = append(statuses, err.Error())
        errors = true
    } else {
        statuses = append(statuses, status)
    }

    status, err = uploadData(cluster, project, "prproj", prproj)
    if err != nil {
        statuses = append(statuses, err.Error())
        errors = true
    } else {
        statuses = append(statuses, status)
    }

    statusString = strings.Join(statuses, ",")

    return
}

func (p Project)NormalisedName()(string) {
    return strings.Replace(p.Name, " ", "_", -1)
}
