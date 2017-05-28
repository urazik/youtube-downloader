package download

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	u "net/url"
	"os"
	"path/filepath"
	"strings"
)

// errors const
const (
	fpErr = "File parsing error"
	upErr = "Url parsing error"
	cfErr = "Create file error"
	dErr  = "Donwload error"
)

// for progress bar
type progress struct {
	io.Reader
	current int64
	total   int64
	title   string
}

// update progress bar
func (p *progress) Read(b []byte) (int, error) {
	var bytesToMB float64 = 1048576

	n, err := p.Reader.Read(b)
	if err == nil {
		p.current += int64(n)

		fmt.Printf(
			"%s   %vMB / %vMB   %v\r",
			p.title,
			fmt.Sprintf("%.2f", float64(p.current)/bytesToMB),
			fmt.Sprintf("%.2f", float64(p.total)/bytesToMB),
			int(float64(p.current)/float64(p.total)*float64(100)+1),
		)
	}

	return n, err
}

// get urls from file
func getFileUrls(path string) []string {
	// get abs path
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(fpErr)
	}

	// read file
	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal(fpErr)
	}

	// parse urls
	urls := make([]string, 0)
	for _, s := range strings.Split(string(file), "\n") {
		urls = append(urls, s)
	}

	return urls
}

// parse one url
func urlParsing(url string) (string, string) {
	// parse youtube video url
	vUrl, err := u.Parse(url)
	if err != nil {
		log.Fatal(upErr)
	}
	m, err := u.ParseQuery(vUrl.RawQuery)
	if err != nil {
		log.Fatal(upErr)
	}
	if len(m) == 0 {
		log.Fatal(upErr)
	}

	// get video info
	resp, err := http.Get("http://www.youtube.com/get_video_info?video_id=" + m["v"][0])
	if err != nil {
		log.Fatal(upErr)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(upErr)
	}

	// parse video meta
	meta, err := u.ParseQuery(string(body))
	if err != nil {
		log.Fatal(upErr)
	}
	if meta["status"][0] == "fail" {
		log.Fatal(upErr)
	}

	// get download video urls
	dUrls, err := u.ParseQuery(meta["url_encoded_fmt_stream_map"][0])
	if err != nil {
		log.Fatal(upErr)
	}

	// to create a file in windows
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

	return dUrls["url"][0], r.Replace(meta["title"][0])
}

// create file
func createFile(dir, title string) *os.File {
	// check for path separator
	if strings.LastIndex(dir, string(filepath.Separator)) != len(dir)-1 {
		dir = dir + string(filepath.Separator)
	}

	// check for file exists
	path := dir + title + ".mp4"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// if exists -> remove
		err := os.Remove(path)
		if err != nil {
			log.Fatal(cfErr)
		}
	}

	// create file
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(cfErr)
	}
	return out
}

// Start start videos download
func Start(url, file, dir string) {
	var urls []string

	// cli returns url or file
	if url != "" {
		urls = append(urls, url)
	} else {
		urls = getFileUrls(file)
	}

	fmt.Printf("Start downloading %v videos", len(urls))

	for _, url := range urls {
		fmt.Println()

		// parse url
		dUrl, title := urlParsing(url)

		// create file
		out := createFile(dir, title)
		defer out.Close()

		// download video
		resp, err := http.Get(dUrl)
		if err != nil {
			log.Fatal(dErr)
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, &progress{
			Reader: resp.Body,
			title:  title,
			total:  resp.ContentLength,
		})
		if err != nil {
			log.Fatal(dErr)
		}
	}

	fmt.Print("\nDownload finished")
}
