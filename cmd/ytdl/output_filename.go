package main

import (
	"bytes"
	"github.com/otium/ytdl"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var funcMap template.FuncMap

func init() {
	funcMap = make(template.FuncMap)
	r := reflect.Indirect(reflect.ValueOf(outputFileName{})).Type()
	for i := 0; i < r.NumField(); i++ {
		v := r.Field(i)
		if v.Type.Kind() == reflect.String {
			funcMap[v.Name] = sanitizeFileNameComponent
		}
	}
}

type outputFileName struct {
	Title            string
	Author           string
	Ext              string
	RawDatePublished time.Time
	Resolution       string
	RawDuration      time.Duration
}

func newOutputFileName(info *ytdl.VideoInfo, format ytdl.Format) *outputFileName {
	return &outputFileName{
		Title:            info.Title,
		Author:           info.Author,
		Ext:              format.Extension,
		RawDatePublished: info.DatePublished,
		Resolution:       format.Resolution,
		RawDuration:      info.Duration,
	}
}

func (f *outputFileName) Duration() string {
	return f.RawDuration.String()
}

func (f *outputFileName) DatePublished() string {
	return f.RawDatePublished.Format("Jan 2 2006")
}

var fileNameTemplate = template.New("OutputFileName")

func (f *outputFileName) String(template string) (string, error) {
	t, err := fileNameTemplate.Parse(template)
	t.Funcs(funcMap)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, f)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var illegalFileNameCharacters = regexp.MustCompile(`[^[a-zA-Z0-9]-_]`)

func sanitizeFileNameComponent(part string) string {
	part = strings.Replace(part, string(os.PathSeparator), "-", -1)
	part = illegalFileNameCharacters.ReplaceAllString(part, "")
	return part
}
