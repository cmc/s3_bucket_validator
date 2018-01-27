/* 
   TODO:
   1. If list successful with invalid keys, log an alert to KAFKA (currently its just writing to a file)
   2. Integrate goroutines & channels to perform all s3 checks concurrently. save time, profit more.
*/
package main

import (
    "log"
    "fmt"
    "regexp"
    "strings"
    "os"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/fatih/color"
)

// return list of all buckets with the key that we're using
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

// return true/false on bucket access status
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
         write_logfile(resp)
        fmt.Println(resp)
    }
    return resp, err

}
func write_logfile(l *s3.ListObjectsOutput) {
    f, err := os.OpenFile("s3_results.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
	fmt.Println("error opening file: %v", err)
    }
    defer f.Close()

    log.SetOutput(f)
    log.Println(l)
}

func main() {
    buckets := get_buckets()
    green := color.New(color.FgGreen).PrintfFunc()

    for _, bucket := range buckets  {
        green("[+] Attempting to access bucket: %s\n", bucket)
        test_bucket_access(bucket, "us-east-1")
        resp, err := test_bucket_access(bucket, "us-east-1")
        if err != nil {
          //  fmt.Println(err)
        }
        if resp != nil {
//            fmt.Println(resp)
        }
    }
}