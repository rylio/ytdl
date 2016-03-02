package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/otium/ytdl"
)

func main() {
	app := cli.NewApp()
	app.Name = "ytdl"
	app.HelpName = "ytdl"
	// Set our own custom args usage
	app.ArgsUsage = "[youtube url or video id]"
	app.Usage = "Download youtube videos"
	app.HideHelp = true
	app.Version = "0.5.0"

	opt := options{}

	app.Flags = []cli.Flag{
		cli.HelpFlag,
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "Write output to a file, passing - outputs to stdout",
			Value:       "{{.Title}}.{{.Ext}}",
			Destination: &opt.Output,
		},
		cli.BoolFlag{
			Name:        "info, i",
			Usage:       "Only output video info",
			Destination: &opt.Info,
		},
		cli.BoolFlag{
			Name:        "no-progress",
			Usage:       "Disable the progress bar",
			Destination: &opt.NoProgress,
		},
		cli.BoolFlag{
			Name:        "silent, s",
			Usage:       "Only output errors, also disables progress bar",
			Destination: &opt.SilentMode,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Output debug log",
			Destination: &opt.DebugMode,
		},
		cli.BoolFlag{
			Name:        "append, a",
			Usage:       "Append to output file instead of overwriting",
			Destination: &opt.Append,
		},
		cli.StringSliceFlag{
			Name:  "filter, f",
			Usage: "Filter available formats, syntax: [format_key]:val1,val2",
		},
		cli.StringFlag{
			Name:        "range, r",
			Usage:       "Download a specific range of bytes of the video, [start]-[end]",
			Destination: &opt.Range,
		},
		cli.BoolFlag{
			Name:        "download-url, u",
			Usage:       "Prints download url to stdout",
			Destination: &opt.DownloadURL,
		},
		cli.BoolFlag{
			Name:        "json, j",
			Usage:       "Print info json to stdout",
			Destination: &opt.JSON,
		},
		cli.StringFlag{
			Name:        "start-offset",
			Usage:       "Offset the start of the video by time",
			Destination: &opt.Offset,
		},
	}

	app.Action = func(c *cli.Context) {
		identifiers := c.Args()
		if len(identifiers) == 0 || c.Bool("help") {
			cli.ShowAppHelp(c)
		} else {
			run(identifiers, opt)
		}
	}
	app.Run(os.Args)
}

func run(identifiers []string, opt options) {

	var logWriter io.Writer = os.Stdout
	if opt.Output == "-" {
		logWriter = os.Stderr
	}
	log.SetOutput(logWriter)

	silentMode := opt.SilentMode || opt.Info || opt.JSON || opt.DownloadURL

	if opt.DebugMode {
		log.SetLevel(log.DebugLevel)
	} else if silentMode {
		log.SetLevel(log.FatalLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	for _, identifier := range identifiers {
		log.Info("Getting video info for ", identifier, "...")
		info, err := ytdl.GetVideoInfo(identifier)
		if err != nil {
			log.Error("Unable to get video info: ", err.Error())
			continue
		}
		if opt.Info || opt.JSON {
			printInfo(info, opt.JSON)
			continue
		}
		formats := info.Formats

		for _, f := range opt.Filters {
			filter, err := parseFilter(f)
			if err == nil {
				formats = filter(formats)
			} else {
				log.Debug("Error parsing format:", err)
			}
		}
		if len(formats) == 0 {
			log.Error("No available formats")
			continue
		}
		selectedFormat := formats[0]
		downloadURL, err := info.GetDownloadURL(selectedFormat)
		if err != nil {
			log.Error("Error getting download url: ", err.Error())
			continue
		}
		if opt.Offset != "" {
			if offset, err := time.ParseDuration(opt.Offset); err != nil {
				query := downloadURL.Query()
				query.Set("begin", fmt.Sprint(int64(offset/time.Millisecond)))
				downloadURL.RawQuery = query.Encode()
			} else {
				log.Debug("Error parsing offset: ", err)
			}
		}
		if opt.DownloadURL {
			fmt.Println(downloadURL.String())
			//	if isatty.IsTerminal(os.Stdout.Fd()) {
			//		fmt.Println()
			//	}
			continue
		}
		var outputWriter io.Writer
		if opt.Output == "-" {
			outputWriter = os.Stdout
		} else {
			fileName, err := newOutputFileName(info, selectedFormat).String(opt.Output)
			if err != nil {
				log.Error("Error creating output filename: ", err.Error())
				continue
			}
			flags := os.O_CREATE | os.O_WRONLY
			if opt.Append {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}
			outputWriter, err = os.OpenFile(fileName, flags, 0666)
			if err != nil {
				log.Error("Error opening output file: ", err.Error())
				continue
			}
			defer outputWriter.(*os.File).Close()
		}
		req, err := http.NewRequest("GET", downloadURL.String(), nil)
		if err != nil {
			log.Error("Error creating request: ", err)
			continue
		}
		if opt.Range == "" && opt.Append {
			if outputWriter != os.Stdout {
				stat, err := outputWriter.(*os.File).Stat()
				if err == nil {
					opt.Range = strconv.FormatInt(stat.Size(), 10) + "-"
				}
			}
		}
		if opt.Range != "" {
			req.Header.Set("Range", "bytes="+opt.Range)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("Error starting download: ", err.Error())
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Error("Recieved invalid status code: ", resp.StatusCode)
			continue
		}
		var progressBar *pb.ProgressBar
		if !silentMode && !opt.NoProgress {
			progressBar = pb.New64(resp.ContentLength).SetUnits(pb.U_BYTES)
			progressBar.ShowTimeLeft = true
			progressBar.ShowSpeed = true
			progressBar.Start()
			outputWriter = io.MultiWriter(outputWriter, progressBar)
		}
		_, err = io.Copy(outputWriter, resp.Body)
		if progressBar != nil {
			progressBar.Finish()
		}
		if err != nil {
			log.Error("Error downloading video: ", err)
		}
	}
}

func printInfo(info *ytdl.VideoInfo, json bool) {
	fmt.Println(info.Title)
}
