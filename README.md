# s3_bucket_validator
Obtains list of available s3 buckets with an internal account, then tries to access each bucket (ListObjects) with an external account. Writes JSON of bucket name + contents returned to KAFKA for alerting.

Internal account used should at a minimum have access to list buckets.

External account used should NOT have access to the buckets in the list if they are locked down, thus, successful accesses by the external account can be used to generate alerts.

Uses AWS Golang SDK

Requires ~/.aws/credentials config with sections:

    [internal-privileged-account]

    [external-unprivileged-account]


To run in container with ~/.aws/credentials configured on local machine:

    $ docker build -t s3test .

    $ docker run -it --rm -e "HOME=/home" -v $HOME/.aws:/home/.aws --name s3tester s3tester
