# ytdl
------
Go library for downloading YouTube videos

[Documentation: https://godoc.org/github.com/otium/ytdl](https://godoc.org/github.com/otium/ytdl "ytdl")

## ytdl CLI

To install: ``go get -u github.com/otium/ytdl/...``
```
NAME:
   ytdl - Download youtube videos

USAGE:
   ytdl [global options] [youtube url or video id]

VERSION:
   0.0.1

GLOBAL OPTIONS:
   --help, -h						show help
   --output, -o "{{.Title}}.{{.Ext}}"			Write output to a file, passing - outputs to stdout
   --info, -i						Only output video info
   --no-progress					Disable the progress bar
   --silent, -s						Only output errors, also disables progress bar
   --debug, -d						Output debug log
   --append, -a						Append to output file instead of overwriting
   --filter, -f [--filter option --filter option]	Filter available formats, syntax: [format_key]:val1,val2
   --range, -r 						Download a specific range of bytes of the video, [start]-[end]
   --download-url, -u					Prints download url to stdout
   --version, -v					print the version
```

## Example
```
import (
   "github.com/otium/ytdl"
   "os"
)

info, err := ytdl.GetInfo("https://www.youtube.com/watch?v=1rZ-JorHJEY")
file, _ = os.Create(info.Title + ".mp4")
defer file.Close()
info.Download(file)

```

## License
ytdl is released under the MIT License, see LICENSE for more details.
