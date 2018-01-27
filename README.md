# s3_bucket_validator
Obtains buckets with an internal account, tests access with an external account.

Uses AWS Golang SDK + aws credentials config file

Requires aws credentials config file sections:

[internal-privileged-account]
[external-unprivileged-account]

Writes JSON of bucket name + contents returned to KAFKA

$ docker build -t s3test .
$ docker run -it --rm --name s3test s3test
