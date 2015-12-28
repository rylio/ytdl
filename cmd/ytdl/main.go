package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/mattn/go-isatty"
	"github.com/otium/ytdl"
)

type options struct {
	noProgress  bool
	outputFile  string
	infoOnly    bool
	silent      bool
	debug       bool
	append      bool
	filters     []string
	downloadURL bool
	byteRange   string
	json        bool
}

func main() {
	app := cli.NewApp()
	app.Name = "ytdl"
	app.HelpName = "ytdl"
	// Set our own custom args usage
	app.ArgsUsage = "[youtube url or video id]"
	app.Usage = "Download youtube videos"
	app.HideHelp = true
	app.Version = "0.4.0"

	app.Flags = []cli.Flag{
		cli.HelpFlag,
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Write output to a file, passing - outputs to stdout",
			Value: "{{.Title}}.{{.Ext}}",
		},
		cli.BoolFlag{
			Name:  "info, i",
			Usage: "Only output video info",
		},
		cli.BoolFlag{
			Name:  "no-progress",
			Usage: "Disable the progress bar",
		},
		cli.BoolFlag{
			Name:  "silent, s",
			Usage: "Only output errors, also disables progress bar",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Output debug log",
		},
		cli.BoolFlag{
			Name:  "append, a",
			Usage: "Append to output file instead of overwriting",
		},
		cli.StringSliceFlag{
			Name:  "filter, f",
			Usage: "Filter available formats, syntax: [format_key]:val1,val2",
		},
		cli.StringFlag{
			Name:  "range, r",
			Usage: "Download a specific range of bytes of the video, [start]-[end]",
		},
		cli.BoolFlag{
			Name:  "download-url, u",
			Usage: "Prints download url to stdout",
		},
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "Print info json to stdout",
		},
	}

	app.Action = func(c *cli.Context) {
		identifier := c.Args().First()
		if identifier == "" || c.Bool("help") {
			cli.ShowAppHelp(c)
		} else {
			options := options{
				noProgress:  c.Bool("no-progress"),
				outputFile:  c.String("output"),
				infoOnly:    c.Bool("info"),
				silent:      c.Bool("silent"),
				debug:       c.Bool("debug"),
				append:      c.Bool("append"),
				filters:     c.StringSlice("filter"),
				downloadURL: c.Bool("download-url"),
				byteRange:   c.String("range"),
				json:        c.Bool("json"),
			}
			if len(options.filters) == 0 {
				options.filters = cli.StringSlice{
					fmt.Sprintf("%s:mp4", ytdl.FormatExtensionKey),
					fmt.Sprintf("%s:1080p,720p,480p,360p,240p,144p", ytdl.FormatResolutionKey),
					fmt.Sprintf("!%s:nil", ytdl.FormatVideoEncodingKey),
					fmt.Sprintf("!%s:nil", ytdl.FormatAudioEncodingKey),
				}
			}
			handler(identifier, options)
		}
	}
	app.Run(os.Args)
}

func handler(identifier string, options options) {
	var err error
	defer func() {
		if err != nil {
			log.SetOutput(os.Stderr)
			log.Fatal(err.Error())
		}
	}()

	var out io.Writer
	var logOut io.Writer = os.Stdout
	// if downloading to stdout, set log output to stderr, not sure if this is correct
	if options.outputFile == "-" {
		out = os.Stdout
		logOut = os.Stderr
	}
	log.SetOutput(logOut)

	// ouput only errors or not
	silent := options.outputFile == "" ||
		options.silent || options.infoOnly || options.downloadURL || options.json
	if silent {
		log.SetLevel(log.FatalLevel)
	} else if options.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// TODO: Show activity indicator
	log.Info("Fetching video info...")
	//fmt.Print("\u001b[0G")
	//fmt.Print("\u001b[2K")
	info, err := ytdl.GetVideoInfo(identifier)
	if err != nil {
		err = fmt.Errorf("Unable to fetch video info: %s", err.Error())
		return
	}

	// TODO: Get more info, and change structure
	// TODO: Allow json output
	if options.infoOnly {
		fmt.Println("Title:", info.Title)
		fmt.Println("Author:", info.Author)
		fmt.Println("Date Published:", info.DatePublished.Format("Jan 2 2006"))
		fmt.Println("Duration:", info.Duration)
		return
	} else if options.json {
		var data []byte
		data, err = json.MarshalIndent(info, "", "\t")
		if err != nil {
			err = fmt.Errorf("Unable to marshal json: %s", err.Error())
			return
		}
		fmt.Println(string(data))
		return
	}

	formats := info.Formats
	// parse filter arguments, and filter through formats
	for _, filter := range options.filters {
		key, values, exclude, err := parseFilter(filter)
		if err == nil {
			if exclude {
				formats = ytdl.FilterFormatsExclude(formats, ytdl.FormatKey(key), values)
			} else {
				formats = ytdl.FilterFormats(formats, ytdl.FormatKey(key), values)
			}
		}
	}
	if len(formats) == 0 {
		err = fmt.Errorf("No formats available that match criteria")
		return
	}
	format := formats[0]
	downloadURL, err := info.GetDownloadURL(format)
	if err != nil {
		err = fmt.Errorf("Unable to get download url: %s", err.Error())
		return
	}
	if options.downloadURL {
		fmt.Print(downloadURL.String())
		// print new line character if outputing to terminal
		if isatty.IsTerminal(os.Stdout.Fd()) {
			fmt.Println()
		}
		return
	}
	if out == nil {
		var fileName string
		fileName, err = createFileName(options.outputFile, outputFileName{
			Title:         info.Title,
			Ext:           format.Extension,
			DatePublished: info.DatePublished.Format("2006-01-02"),
			Resolution:    format.Resolution,
			Author:        info.Author,
			Duration:      info.Duration.String(),
		})
		if err != nil {
			err = fmt.Errorf("Unable to parse output file file name: %s", err.Error())
			return
		}
		// Create file truncate if append flag is not set
		flags := os.O_CREATE | os.O_WRONLY
		if options.append {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}
		var f *os.File
		// open as write only
		f, err = os.OpenFile(fileName, flags, 0666)
		if err != nil {
			err = fmt.Errorf("Unable to open output file: %s", err.Error())
			return
		}
		defer f.Close()
		out = f
	}

	log.Info("Downloading to ", out.(*os.File).Name())
	var req *http.Request
	req, err = http.NewRequest("GET", downloadURL.String(), nil)
	// if byte range flag is set, use http range header option
	if options.byteRange != "" || options.append {
		if options.byteRange == "" && out != os.Stdout {
			if stat, err := out.(*os.File).Stat(); err != nil {
				options.byteRange = strconv.FormatInt(stat.Size(), 10) + "-"
			}
		}
		if options.byteRange != "" {
			req.Header.Set("Range", "bytes="+options.byteRange)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if err == nil {
			err = fmt.Errorf("Received status code %d from download url", resp.StatusCode)
		}
		err = fmt.Errorf("Unable to start download: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	// if we aren't in silent mode or the no progress flag wasn't set,
	// initialize progress bar
	if !silent && !options.noProgress {
		progressBar := pb.New64(resp.ContentLength)
		progressBar.SetUnits(pb.U_BYTES)
		progressBar.ShowTimeLeft = true
		progressBar.ShowSpeed = true
		//	progressBar.RefreshRate = time.Millisecond * 1
		progressBar.Output = logOut
		progressBar.Start()
		defer progressBar.Finish()
		out = io.MultiWriter(out, progressBar)
	}
	_, err = io.Copy(out, resp.Body)
}
