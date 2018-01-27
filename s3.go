package main

import (
    "os"
    "fmt"
    "strings"
    "regexp"
    "encoding/json"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/fatih/color"
)

// return list of all buckets with the privileged internal account
func get_buckets() []string {
    os.Setenv("AWS_PROFILE", "internal-privileged-account")
    s3svc := s3.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
    result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
    if err != nil {
            fmt.Println("Failed to list buckets", err)
    }
    bucket_list := make([]string, len(result.Buckets))
    for i, bucket := range result.Buckets {
         bucket_list[i] = aws.StringValue(bucket.Name)
    }
    fmt.Println("Retrieved all buckets: ", bucket_list)
    return bucket_list
}

// test whether we can list bucket objects with the unprivileged, external account
func test_bucket_access(bucket string, region string) (*s3.ListObjectsOutput, error) {
    red := color.New(color.FgRed).PrintfFunc()
    yellow := color.New(color.FgYellow).PrintfFunc()

    os.Setenv("AWS_PROFILE", "external-unprivileged-account")
    svc := s3.New(session.New(), &aws.Config{Region: aws.String(region)})
    params := &s3.ListObjectsInput{
          Bucket: aws.String(bucket),
    }
    resp, err := svc.ListObjects(params)
    if err != nil {
        if strings.Contains(err.Error(), "is wrong") {
            re := regexp.MustCompile("'([^' ]+)'")
            correct_region := re.FindAllStringSubmatch(err.Error(), 2)[1][1]
            yellow("[-] Wrong region!, retrying with region: %s\n", correct_region)
            test_bucket_access(bucket, correct_region)
        } else {
            yellow("[-] failed to access bucket: %s with error: %s\n", bucket, err)
        }
    }
    if err == nil && resp != nil {
        red("[!] WARNING, Access Granted to bucket: %s\n", bucket)
        mapD := map[string]interface{}{"type": "alert", "bucket": bucket, "raw_data": resp}
        mapB, _ := json.Marshal(mapD)
        if send_kafka(mapB) {
            fmt.Println("Successfully written JSON alert to kafka.")
        }
        write_logfile(resp)
        fmt.Println(resp)
    }
    return resp, err

}
