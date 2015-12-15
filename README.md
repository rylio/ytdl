
# ytdl [![Build Status](https://travis-ci.org/otium/ytdl.svg)](https://travis-ci.org/otium/ytdl) [![GoDoc](https://godoc.org/github.com/codegangsta/cli?status.svg)](https://godoc.org/github.com/otium/ytdl)
------
Go library for downloading YouTube videos

[Documentation: https://godoc.org/github.com/otium/ytdl](https://godoc.org/github.com/otium/ytdl "ytdl")

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

## ytdl CLI

To install: ``go get -u github.com/otium/ytdl/...``

### Usage
``` ytdl [global options] [youtube url or video id] ```
### Options
 - ```--help, -h``` - show help
 - ```--filter, -f``` - Filter out formats
   - Syntax: key:value1,value2,...,valueN
   - To exclude: !key:value1,...
   - Available keys (See format.go for available values):
      - ```ext``` - extension of video
      - ```res``` - resolution of video
      - ```videnc``` - video encoding
      - ```audenc``` - audio encoding
      - ```prof``` - youtube video profile
   - Default filters
      - ```ext:mp4```
      - ```res:1080p,720p,480p,360p,240p,144p```
      - ```!videnc:nil```
      - ```!audenc:nil```
 - ```--output, -o``` - Output to specific path
   - Supports templates, ex: {{.Title}}.{{.Ext}}
   - Defaults to ```{{.Title}}.{{.Ext}}```
   - Supported template variables are Title, Ext, DatePublished, Resolution
   - Pass - to output to stdout, former stdout output is redirected to stderr
 - ```--info, -i``` - Just gets video info, outputs to stdout
 - ```--no-progress``` - Disables the progress bar
 - ```--silent, -s``` - Disables all output, except for fatal errors
 - ```--debug, -d``` - Output debug logs
 - ```--append, -a``` - append to output file, instead of truncating
 - ```--range, -r``` - specify a range of bytes, placed in http range header, ex: 0-100
 - ```--download-url, -u``` - just print download url to, don't do anything else
 - ```--version, -v``` - print out ytdl cli version

## License
ytdl is released under the MIT License, see LICENSE for more details.
