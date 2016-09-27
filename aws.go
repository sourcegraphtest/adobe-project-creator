package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "bytes"
    "fmt"
)


func uploadData(c Cluster, p Project, fileSuffix string, projectData []byte)(string, error) {
    svc := s3.New(session.New(), &aws.Config{Region: aws.String("eu-west-1")})

    params := &s3.PutObjectInput{
        Bucket:             aws.String(c.DestBucket()),
        Key:                aws.String( fmt.Sprintf("%s/%s.%s", p.UUID, p.NormalisedName(), fileSuffix)),
        Body:               bytes.NewReader(projectData),
    }

    output, err := svc.PutObject(params)
    return output.String(), err
}
