package main

import (
	"fmt"
	"github.com/matthew-andrews/s3up/objects"
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
			Name:  "region",
			Usage: "Optionally set the region (defaults to whatever set by the environment)",
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
		if len(files) < 1 {
			return cli.NewExitError("No files found for upload to S3.  (Directories are ignored)", 1)
		}
		for _, file := range files {
			fmt.Printf("%s to %s\n", file.Location, file.Key)
		}
		return nil
	}
	app.Run(os.Args)
}
