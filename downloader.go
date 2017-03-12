package main

import (
	"log"
	"net/url"
	"io/ioutil"
	"net/http"
	"os"
	"io"
)

func main() {
	videoURL := "https://www.youtube.com/watch?v=cUBMQznYuBM"
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

	out, err := os.Create("video.mp4")
	check(err)
	defer out.Close()

	resp, err = http.Get(u3["url"][0])
	check(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
