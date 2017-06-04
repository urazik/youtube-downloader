# YouTube video downloader
This application allows you to download videos from YouTube

## Usage
### CLI
#### Download by url
```
$ go run cli.go -url=https://www.youtube.com/watch?v=<video_id>
```
#### Download from file
```
$ go run cli.go -file=<path>
```
##### Example of a file structure
```
https://www.youtube.com/watch?v=<video_id>
https://www.youtube.com/watch?v=<video_id>
https://www.youtube.com/watch?v=<video_id>
```
#### Download to a specific directory
```
$ go run cli.go -url=https://www.youtube.com/watch?v=<video_id> -dir=<path>
```
### Go
```go
package main

import (
	"github.com/lavrs/youtube-downloader"
)

func main() {
	var url string = "https://www.youtube.com/watch?v=<video_id>"
	var dir string = "/home/lavrs/download"
	download.Start(url, "", dir)
	
	var file = "/home/lavrs/videos.txt"
	download.Start("", file, dir)
}
```