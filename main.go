package main

import (
	"fmt"
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
			Usage: "Optionally set the Canned Access Control List for new files being put into S3 (default to public-read)",
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("will upload some files, I guess")
		return nil
	}
	app.Run(os.Args)
}
