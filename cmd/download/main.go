package main

import (
    d "github.com/lavrs/youtube-downloader"
    "github.com/urfave/cli"
    "os"
)

func main() {
    app := cli.NewApp()
    app.Usage = "dowload video from YouTube"

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "u, url",
            Usage: "YouTube video url",
        },

        cli.StringFlag{
            Name:  "d, dir",
            Usage: "path to download dir",
        },

        cli.StringFlag{
            Name:  "f, file",
            Usage: "file with urls",
        },
    }

    app.Action = func(c *cli.Context) {
        if c.NArg() == 0 {
            var urls []string

            if c.String("u") != "" {
                urls = append(urls, c.String("u"))
            } else if c.String("f") != "" {
                urls = d.GetUrlsFromFile(c.String("f"))
            } else {
                return
            }

            if c.String("d") != "" {
                d.Download(urls, c.String("d"))
                return
            } else {
                d.Download(urls, "")
            }
        } else {
            cli.ShowAppHelp(c)
        }
    }

    app.Run(os.Args)
}
