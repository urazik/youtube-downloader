package main

import (
	"flag"
	"fmt"
	d "github.com/lavrs/youtube-downloader"
)

func main() {
	var (
		url  = flag.String("url", "", "YouTube video url")
		dir  = flag.String("dir", "", "Path to download dir")
		file = flag.String("file", "", "File with urls")
	)
	flag.Parse()

	if *file == "" || *url == "" {
		fmt.Println("Usage error! Insert url or file")
		return
	}
	d.Start(*url, *file, *dir)
}
