package main

import (
	"net/url"
	"io/ioutil"
	"net/http"
	"os"
	"io"
	"fmt"
	"strings"
	"github.com/urfave/cli"
	"errors"
)

type progress struct {
	io.Reader
	current int64
	title   string
	length  int64
}

func (pr *progress) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)

	if err == nil {
		pr.current += int64(n)
		var bytesToMB float64 = 1048576

		fmt.Printf(
			"%s   %vMB / %vMB   %v%%\r",
			pr.title,
			fmt.Sprintf("%.2f", float64(pr.current)/bytesToMB),
			fmt.Sprintf("%.2f", float64(pr.length)/bytesToMB),
			int(float64(pr.current)/float64(pr.length)*float64(100)+1),
		)
	}

	return n, err
}

func main() {
	app := cli.NewApp()
	app.Name = "yvd"
	app.Usage = "dowload video from YouTube"
	app.Version = "0.1.0"

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
				urls = openfile(c.String("f"))
			}

			if c.String("d") != "" {
				download(urls, c.String("d"))
				return
			} else {
				download(urls, "")
			}
		} else {
			cli.ShowAppHelp(c)
		}
	}

	err := app.Run(os.Args)
	check(err)
}

func openfile(path string) []string {
	urls := make([]string, 0)
	file, err := ioutil.ReadFile(path)
	check(err)

	for _, s := range strings.Split(string(file), "\n") {
		urls = append(urls, s)
	}

	return urls
}

func urlParsing(u string) (string, string, error) {
	videoUrl, err := url.Parse(u)
	check(err)
	m, err := url.ParseQuery(videoUrl.RawQuery)
	check(err)

	if len(m) == 0 {
		return "", "", errors.New("parse error")
	}

	getVideoInfoURL := "http://www.youtube.com/get_video_info?video_id=" + m["v"][0]
	resp, err := http.Get(getVideoInfoURL)
	check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	metadata, err := url.ParseQuery(string(body))
	check(err)

	if metadata["status"][0] == "fail" {
		return "", "", errors.New("parse error")
	}

	URLEncodedFmtStreamMap, err := url.ParseQuery(metadata["url_encoded_fmt_stream_map"][0])
	check(err)

	r := strings.NewReplacer(
		"\\", "",
		"/", "",
		":", "",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	title := r.Replace(metadata["title"][0])

	return URLEncodedFmtStreamMap["url"][0], title, nil
}

func createFile(path string) *os.File {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := os.Remove(path)
		check(err)
	}

	out, err := os.Create(path)
	check(err)
	defer out.Close()

	return out
}

func download(urls []string, dir string) {
	for _, u := range urls {
		fmt.Println()

		dUrl, title, err := urlParsing(u)
		check(err)

		out := createFile(dir + title + ".mp4")

		resp, err := http.Get(dUrl)
		check(err)
		defer resp.Body.Close()

		_, err = io.Copy(out, &progress{
			Reader: resp.Body,
			title:  title,
			length: resp.ContentLength,
		})
		check(err)
	}
}

func check(err error) {
	if err != nil {
		fmt.Println("Oops, some error! Check the urls")
		os.Exit(0)
	}
}
