package main

import (
	"fmt"
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
	app.Version = "v1.0.4"
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
	}
	app.Action = func(c *cli.Context) error {
		files, _ := objects.GetFiles(c.Args(), c.Int("strip"), c.String("destination"), c.String("cache-control"), c.String("acl"))
		awsSession := session.New()
		service := s3client.New(s3.New(awsSession), s3manager.NewUploader(awsSession))
		err := uploader.Upload(service, c.String("bucket"), files, c.Int("concurrency"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err), 1)
		}
		return nil
	}
	app.Run(os.Args)
}
