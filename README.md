# YouTube video downloader
This application allows you to download videos from YouTube

## Usage
### CLI
#### Download by url
```
$ go run main.go -u https://www.youtube.com/watch?v=<video_id>
```
#### Download from file
```
$ go run main.go --file <path>
```
##### Example of a file structure
```
https://www.youtube.com/watch?v=<video_id>
https://www.youtube.com/watch?v=<video_id>
https://www.youtube.com/watch?v=<video_id>
...
```
#### Download to a specific directory
```
$ go run main.go --url https://www.youtube.com/watch?v=<video_id> --dir <path>
```
### Go
```go
package main

import (
	"github.com/lavrs/youtube-downloader"
)

func main() {
	var urls []string = []string{"https://www.youtube.com/watch?v=<video_id>"}
	var dir string = "/home/lavrs/download"

	download.Download(urls, dir)
}
```