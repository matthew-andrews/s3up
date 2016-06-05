package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"github.com/matthew-andrews/s3up/s3client"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "S3 Upload"
	app.Usage = "Uploads files to S3"
	app.Version = "v1.0.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "strip",
			Usage: "Optionally remove the specified number of leading path elements",
		},
		cli.StringFlag{
			Name:  "destination",
			Usage: "Optionally add a prefix to the upload path",
		},
		cli.StringFlag{
			Name:  "bucket",
			Usage: "Specify the S3 bucket to upload files to",
		},
		cli.StringFlag{
			Name:  "cache-control",
			Usage: "Optionally set a Cache-Control value",
		},
		cli.StringFlag{
			Name:  "acl",
			Value: "public-read",
			Usage: "Optionally set the Canned Access Control List for new files being put into S3 (default to public-read)",
		},
	}
	app.Action = func(c *cli.Context) error {
		files, _ := objects.GetFiles(c.Args(), c.Int("strip"), c.String("destination"), c.String("cache-control"), c.String("acl"))
		service := s3client.Client{
			Service: s3.New(session.New()),
		}
		err := service.Upload(c.String("bucket"), files)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err), 1)
		}
		return nil
	}
	app.Run(os.Args)
}
