package download

import (
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "log"
)

type progress struct {
    io.Reader
    current int64
    title   string
    length  int64
}

func (pr *progress) Read(p []byte) (int, error) {
    var (
        n   int
        err error
    )

    if n, err = pr.Reader.Read(p); err == nil {
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

func GetUrlsFromFile(path string) []string {
    absPath, err := filepath.Abs(path)
    check(err)
    file, err := ioutil.ReadFile(absPath)
    check(err)

    urls := make([]string, 0)

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

func createFile(dir, title string) *os.File {
    if strings.LastIndex(dir, string(filepath.Separator)) != len(dir)-1 {
        dir = dir + string(filepath.Separator)
    }

    path := dir + title + ".mp4"

    if _, err := os.Stat(path); !os.IsNotExist(err) {
        err := os.Remove(path)
        check(err)
    }

    out, err := os.Create(path)
    check(err)

    return out
}

func Download(urls []string, dir string) {
    fmt.Print("Start download\n")

    for _, u := range urls {
        fmt.Println()

        dUrl, title, err := urlParsing(u)
        check(err)

        out := createFile(dir, title)
        defer out.Close()

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

    fmt.Print("\n\nDownload end")
}

func check(err error) {
    if err != nil {
        log.Fatal("Oops, some error! Check the urls or file path.", err)
    }
}
