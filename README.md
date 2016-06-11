# s3up

## Installation

Grab a binary from:- https://github.com/matthew-andrews/s3up/releases and put it on your `PATH`.

## Usage

```
NAME:
   S3 Upload - Uploads files to S3

USAGE:
   s3up [global options] command [command options] [arguments...]
   
VERSION:
   v1.0.4

COMMANDS:
GLOBAL OPTIONS:
   --strip value          optionally remove the specified number of leading path elements (default: 0)
   --concurrency value    optionally configure the maximum number of simultaneous uploads (default: 10)
   --destination value    optionally add a prefix to the upload path
   --bucket value         specify the S3 bucket to upload files to
   --cache-control value  optionally set a Cache-Control value
   --acl value            optionally set the Canned Access Control List for new files being put into S3 (default to public-read) (default: "public-read")
   --help, -h             show help
   --version, -v          print the version
``` 

## Examples

The following command will recursively find all the files in the `public` subfolder, strip off the `./public` prefix, add a prefix of `hello/`, set a very long `cache-control` header.

```
export AWS_REGION="eu-west-1"; \
	export AWS_ACCESS_KEY_ID="???"; \
	export AWS_SECRET_ACCESS_KEY="???"; \
	s3up --destination hello --bucket s3up-test --cache-control "public, max-age=31536000" --strip 1 \
	`find . -path "./public/*"`
```

So, if there were a file called:-

```
./public/sub-folder/hello-world.txt`
```

This will be uploaded to:-

```
https://s3-eu-west-1.amazonaws.com/s3up-test/hello/sub-folder/hello-world.txt
```
