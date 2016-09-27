package main

import (
    "fmt"
)

type Cluster struct {
    ID string
}

func (c Cluster) DestBucket() string {
    return fmt.Sprintf("mio-project-%s", c.ID)
}
