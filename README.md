# s3up

## Installation

Grab a binary from:- https://github.com/matthew-andrews/s3up/releases and put it on your `PATH`.

For example, on OS X you can do this via:-

```sh
curl -sfL https://github.com/matthew-andrews/s3up/releases/download/v1.0.5/s3up_darwin_386 -o /usr/local/bin/s3up && chmod +x /usr/local/bin/s3up
```

Or on Ubuntu (or in CI environments such as Circle CI) you can do this via:-

```sh
curl -sfL https://github.com/matthew-andrews/s3up/releases/download/v1.0.5/s3up_linux_386 -o /home/ubuntu/bin/s3up && chmod +x /home/ubuntu/bin/s3up
```

## Usage

```
NAME:
   S3 Upload - Uploads files to S3

USAGE:
   s3up [global options] command [command options] [arguments...]
   
VERSION:
   v1.0.5

COMMANDS:
GLOBAL OPTIONS:
   --strip value          optionally remove the specified number of leading path elements (default: 0)
   --concurrency value    optionally configure the maximum number of simultaneous uploads (default: 10)
   --destination value    optionally add a prefix to the upload path
   --bucket value         specify the S3 bucket to upload files to
   --cache-control value  optionally set a Cache-Control value
   --acl value            optionally set the Canned Access Control List for new files being put into S3 (default to public-read) (default: "public-read")
   --dry-run              perform a trial run with no changes made
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

## IAM permissions

s3up integrates with all the usual ways of configuring AWS permissions (`.aws/credentials` in your home directory, `AWS_ACCESS_KEY_ID` / `AWS_SECRET_KEY_ID` environment variables, etc).  The user or role that s3up is authenticating with will also need permissions to read and write files and metadata about files in the buckets it is intended to be used with.
