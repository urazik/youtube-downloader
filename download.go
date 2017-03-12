package main

import (
	"log"
	"net/url"
	"io/ioutil"
	"net/http"
	"os"
	"io"
	"fmt"
	"strings"
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
	video := [...]string{"https://www.youtube.com/watch?v=cUBMQznYuBM", "https://www.youtube.com/watch?v=z9vUCXGoC6w"}

	for _, v := range video {
		fmt.Println()

		u, err := url.Parse(v)
		check(err)
		m, err := url.ParseQuery(u.RawQuery)
		check(err)

		getVideoInfoURL := "https://www.youtube.com/get_video_info?video_id=" + m["v"][0]
		resp, err := http.Get(getVideoInfoURL)
		check(err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		metadata, err := url.ParseQuery(string(body))
		check(err)

		URLEncodedFmtStreamMap, err := url.ParseQuery(metadata["url_encoded_fmt_stream_map"][0])
		check(err)

		title := strings.Replace(metadata["title"][0], ":", "", -1)

		if _, err := os.Stat(title + ".mp4"); !os.IsNotExist(err) {
			err := os.Remove(title + ".mp4")
			check(err)
		}

		out, err := os.Create(title + ".mp4")
		check(err)
		defer out.Close()

		resp, err = http.Get(URLEncodedFmtStreamMap["url"][0])
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
		log.Panic(err)
	}
}
