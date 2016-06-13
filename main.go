package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/matthew-andrews/s3up/objects"
	"github.com/matthew-andrews/s3up/s3client"
	"github.com/matthew-andrews/s3up/uploader"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "S3 Upload"
	app.Usage = "Uploads files to S3"
	app.Version = "v1.0.5"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "strip",
			Usage: "optionally remove the specified number of leading path elements",
		},
		cli.IntFlag{
			Name:  "concurrency",
			Usage: "optionally configure the maximum number of simultaneous uploads",
			Value: 10,
		},
		cli.StringFlag{
			Name:  "destination",
			Usage: "optionally add a prefix to the upload path",
		},
		cli.StringFlag{
			Name:  "bucket",
			Usage: "specify the S3 bucket to upload files to",
		},
		cli.StringFlag{
			Name:  "cache-control",
			Usage: "optionally set a Cache-Control value",
		},
		cli.StringFlag{
			Name:  "acl",
			Value: "public-read",
			Usage: "optionally set the Canned Access Control List for new files being put into S3 (default to public-read)",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "perform a trial run with no changes made",
		},
	}
	app.Action = func(c *cli.Context) error {
		files, _ := objects.GetFiles(c.Args(), c.Int("strip"), c.String("destination"), c.String("cache-control"), c.String("acl"))

		awsSession := session.New()
		var region string
		if os.Getenv("AWS_REGION") == "" {
			region = "us-east-1"
		} else {
			region = os.Getenv("AWS_REGION")
		}
		s3Service := s3.New(awsSession, &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
			S3UseAccelerate:  aws.Bool(false),
			LogLevel:         aws.LogLevel(aws.LogDebugWithHTTPBody),
			Region:           aws.String(region),
		})

		resp, err := s3Service.GetBucketLocation(&s3.GetBucketLocationInput{
			Bucket: aws.String(c.String("bucket")),
		})
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err), 1)
		}
		fmt.Println(resp)
		return cli.NewExitError("lol k", 1)

		service := s3client.New(s3Service, s3manager.NewUploader(awsSession), c.Bool("dry-run"))
		errs := uploader.Upload(service, c.String("bucket"), files, c.Int("concurrency"))
		if errs != nil {
			return cli.NewExitError(fmt.Sprintf("%s", errs), 1)
		}
		return nil
	}
	app.Run(os.Args)
}
