package main

import (
	"log"
	"net/url"
	"io/ioutil"
	"net/http"
	"os"
	"io"
	"fmt"
)

type progress struct {
	io.Reader
	current int64
	length  int64
}

func (pr *progress) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)

	if err == nil {
		pr.current += int64(n)
		var bytesToMB float64 = 1048576

		fmt.Printf(
			"%vMB / %vMB   %v%%\r",
			fmt.Sprintf("%.2f", float64(pr.current)/bytesToMB),
			fmt.Sprintf("%.2f", float64(pr.length)/bytesToMB),
			int(float64(pr.current)/float64(pr.length)*float64(100)+1),
		)
	}

	return n, err
}

func main() {
	videoURL := "https://www.youtube.com/watch?v=6hseaMlH7RM&t=1s"
	u, err := url.Parse(videoURL)
	check(err)
	m, _ := url.ParseQuery(u.RawQuery)
	getInfoURL := "https://www.youtube.com/get_video_info?video_id=" + m["v"][0]

	resp, err := http.Get(getInfoURL)
	check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	u2, err := url.ParseQuery(string(body))
	check(err)
	u3, err := url.ParseQuery(u2["url_encoded_fmt_stream_map"][0])
	check(err)

	if _, err := os.Stat("video.mp4"); os.IsNotExist(err) {
		err := os.Remove("video.mp4")
		check(err)
	}

	out, err := os.Create("video.mp4")
	check(err)
	defer out.Close()

	resp, err = http.Get(u3["url"][0])
	check(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, &progress{
		Reader: resp.Body,
		length: resp.ContentLength,
	})
	check(err)
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
