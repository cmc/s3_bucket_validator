# s3_bucket_validator
Obtains list of available s3 buckets with an internal account, then tries to access each bucket (ListObjects) with an external account. Writes JSON of bucket name + contents returned to KAFKA for alerting.

Internal account used should have access to list buckets.
External account should NOT have access to buckets if they are locked down as intended, thus, successful accesses are alerts.

Uses AWS Golang SDK

Requires ~/.aws/credentials config with sections:

    [internal-privileged-account]

    [external-unprivileged-account]


To run in container with ~/.aws/credentials configured on local machine:

    $ docker build -t s3test .

    $ docker run -it --rm -e "HOME=/home" -v $HOME/.aws:/home/.aws --name s3tester s3tester
