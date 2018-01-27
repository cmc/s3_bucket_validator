# s3_bucket_validator
Obtains buckets with an internal account, tries to access (ListObjects) with an external account. Writes JSON of bucket name + contents returned to KAFKA for alerting.

Uses AWS Golang SDK

Requires ~/.aws/credentials config with sections:

    [internal-privileged-account]

    [external-unprivileged-account]


To run in container with ~/.aws/credentials configured on local machine:

    $ docker build -t s3test .

    $ docker run -it --rm -e "HOME=/home" -v $HOME/.aws:/home/.aws --name s3tester s3tester
