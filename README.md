# storage-bucket

A very simple storage bucket to be used instead of S3 storage

## Deployment

### With docker

#### Build image

```sh
docker build -t storage-bucket .
```

#### Run image

```sh
docker run -d --name storage-bucket \ 
-v /path/to/bucket:/data \
-p 5000:5000
storage-bucket
```

### Without docker

#### Build

```sh
go build -o storage-bucket ./cmd/main.go
```

### Run

```sh
./storage-bucketsd --pathPrefix path/to/bucket.yaml
```

Note:

- Path prefix should NOT end with a trailing /

The /path/to/bucket should follow the rules below

- The path must be absolute
- Inside the folder should be your bucket.yaml file

## Getting keys

To generate keys send a post request to your endpoint with basic authentication header. This should math the username and password specified in the htpasswd field in bucket.yaml