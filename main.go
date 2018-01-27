// cmc
//TODO:
//  1. Integrate goroutines & channels to perform all s3 checks concurrently. save time, profit more.

package main

import (
    "github.com/fatih/color"
)

func main() {
    buckets := get_buckets()
    green := color.New(color.FgGreen).PrintfFunc()

    for _, bucket := range buckets  {
        green("[+] Attempting to access bucket: %s\n", bucket)
        test_bucket_access(bucket, "us-east-1")
    }
}
