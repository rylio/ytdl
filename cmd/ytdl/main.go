package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/corny/ytdl"
	"github.com/mattn/go-isatty"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

type options struct {
	outputFile     string
	infoOnly       bool
	silent         bool
	debug          bool
	append         bool
	filters        []string
	downloadURL    bool
	byteRange      string
	json           bool
	startOffset    string
	downloadOption bool
}

func main() {
	app := cli.NewApp()
	app.Name = "ytdl"
	app.HelpName = "ytdl"
	// Set our own custom args usage
	app.ArgsUsage = "[youtube url or video id]"
	app.Usage = "Download youtube videos"
	app.HideHelp = true
	app.Version = "0.5.0"

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
		cli.StringFlag{
			Name:  "start-offset",
			Usage: "Offset the start of the video by time",
		},
		cli.BoolFlag{
			Name:  "download-option, p",
			Usage: "Print video and audio download options",
		},
	}

	app.Action = func(c *cli.Context) error {
		identifier := c.Args().First()
		if identifier == "" || c.Bool("help") {
			cli.ShowAppHelp(c)
		} else {
			options := options{
				outputFile:     c.String("output"),
				infoOnly:       c.Bool("info"),
				silent:         c.Bool("silent"),
				debug:          c.Bool("debug"),
				append:         c.Bool("append"),
				filters:        c.StringSlice("filter"),
				downloadURL:    c.Bool("download-url"),
				byteRange:      c.String("range"),
				json:           c.Bool("json"),
				startOffset:    c.String("start-offset"),
				downloadOption: c.Bool("download-option"),
			}
			if len(options.filters) == 0 {
				options.filters = cli.StringSlice{
					fmt.Sprintf("%s:mp4", ytdl.FormatExtensionKey),
					fmt.Sprintf("!%s:", ytdl.FormatVideoEncodingKey),
					fmt.Sprintf("!%s:", ytdl.FormatAudioEncodingKey),
					fmt.Sprint("best"),
				}
			}
			handler(identifier, options)
		}
		return nil
	}
	app.Run(os.Args)
}

func handler(identifier string, options options) {
	var itag int
	var err error
	defer func() {
		if err != nil {
			log.Logger = log.Output(os.Stderr)
			log.Fatal().Err(err).Msg("")
		}
	}()

	var out io.Writer
	var logOut io.Writer = os.Stdout

	// if downloading to stdout, set log output to stderr, not sure if this is correct
	if options.outputFile == "-" {
		out = os.Stdout
		logOut = os.Stderr
	}
	log.Logger = log.Output(logOut)

	// ouput only errors or not
	silent := options.outputFile == "" ||
		options.silent || options.infoOnly || options.downloadURL || options.json
	if silent {
		log.Logger = log.Level(zerolog.FatalLevel)
	} else if options.debug {
		log.Logger = log.Level(zerolog.DebugLevel)
	} else {
		log.Logger = log.Level(zerolog.InfoLevel)
	}

	// TODO: Show activity indicator
	log.Info().Msg("Fetching video info...")
	//fmt.Print("\u001b[0G")
	//fmt.Print("\u001b[2K")
	info, err := ytdl.GetVideoInfo(identifier)
	if err != nil {
		err = fmt.Errorf("Unable to fetch video info: %w", err)
		return
	}

	if options.infoOnly {
		fmt.Println("Title:", info.Title)
		fmt.Println("Uploader:", info.Uploader)
		fmt.Println("Artist:", info.Artist)
		fmt.Println("Song:", info.Song)
		fmt.Println("Date Published:", info.DatePublished.Format("Jan 2 2006"))
		fmt.Println("Duration:", info.Duration)
		return
	} else if options.json {
		var data []byte
		data, err = json.MarshalIndent(info, "", "\t")
		if err != nil {
			err = fmt.Errorf("Unable to marshal json: %w", err)
			return
		}
		fmt.Println(string(data))
		return
	} else if options.downloadOption {
		var data [][]string

		info.Formats.Sort(ytdl.FormatResolutionKey, true)
		for _, format := range info.Formats {
			var fps string
			if format.ValueForKey("fps") == nil {
				fps = "n/a"
			} else {
				fps = format.ValueForKey("fps").(string)
			}

			data = append(data, []string{strconv.Itoa(format.Itag), format.Extension, format.Resolution, fps, format.VideoEncoding, format.AudioEncoding, strconv.Itoa(format.AudioBitrate)})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"itag", "ext", "res", "fps", "vEncoding", "aEncoding", "aBitrate"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output

		fmt.Print("Enter the itag of the file you would like to download or enter 0 to abort: ")
		_, err := fmt.Scanf("%d", &itag)
		if err != nil {
			return
		} else if itag == 0 {
			return
		}
	}

	formats := info.Formats

	if itag != 0 {
		filter, err := parseFilter(fmt.Sprintf("itag:%d", itag))
		if err == nil {
			formats = filter(formats)
		}
	} else {
		// parse filter arguments, and filter through formats
		for _, filter := range options.filters {
			filter, err := parseFilter(filter)
			if err == nil {
				formats = filter(formats)
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
		err = fmt.Errorf("Unable to get download url: %w", err)
		return
	}
	if options.startOffset != "" {
		var offset time.Duration
		offset, err = time.ParseDuration(options.startOffset)
		query := downloadURL.Query()
		query.Set("begin", fmt.Sprint(int64(offset/time.Millisecond)))
		downloadURL.RawQuery = query.Encode()
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
			Title:         sanitizeFileNamePart(info.Title),
			Ext:           sanitizeFileNamePart(format.Extension),
			DatePublished: sanitizeFileNamePart(info.DatePublished.Format("2006-01-02")),
			Resolution:    sanitizeFileNamePart(format.Resolution),
			Duration:      sanitizeFileNamePart(info.Duration.String()),
		})
		if err != nil {
			err = fmt.Errorf("Unable to parse output file file name: %w", err)
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
			err = fmt.Errorf("Unable to open output file: %w", err)
			return
		}
		defer f.Close()
		out = f
	}

	log.Info().Msgf("Downloading to %s", out.(*os.File).Name())
	var req *http.Request
	req, err = http.NewRequest("GET", downloadURL.String(), nil)
	// if byte range flag is set, use http range header option
	if options.byteRange != "" || options.append {
		if options.byteRange == "" && out != os.Stdout {
			if stat, err := out.(*os.File).Stat(); err == nil {
				options.byteRange = strconv.FormatInt(stat.Size(), 10) + "-"
			} else {
				err = fmt.Errorf("Unable to retrieve the existing file's stat: %w", err)
				return
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
		err = fmt.Errorf("Unable to start download: %w", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
}
