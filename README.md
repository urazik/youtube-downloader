# YouTube video downloader
This application allows you to download videos from YouTube

## Usage
### Download from url
```
$ go run download.go --url https://www.youtube.com/watch?v=<video_id>
```
### Download from file
```
$ go run download.go --file <path>
```
#### Example of a file structure
```
https://www.youtube.com/watch?v=<video_id>
https://www.youtube.com/watch?v=<video_id>
https://www.youtube.com/watch?v=<video_id>
...
```
### Download to a specific directory
```
$ go run download.go -url https://www.youtube.com/watch?v=<video_id> --dir <path>
```