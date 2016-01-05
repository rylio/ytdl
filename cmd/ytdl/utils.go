package main

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"text/template"
)

func parseFilter(filter string) (string, []string, bool, error) {
	split := strings.SplitN(filter, ":", 2)
	err := errors.New("Invalid filter")
	if len(split) != 2 {
		return "", nil, false, err
	}
	key := split[0]
	exclude := false
	if key[0] == '!' {
		exclude = true
		key = key[1:]
	}
	values := strings.Split(split[1], ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return key, values, exclude, nil
}

type outputFileName struct {
	Title         string
	Author        string
	Ext           string
	DatePublished string
	Resolution    string
	Duration      string
}

var fileNameTemplate = template.New("OutputFileName")

func createFileName(template string, values outputFileName) (string, error) {
	t, err := fileNameTemplate.Parse(template)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, values)
	if err != nil {
		return "", err
	}
	return string(buf.String()), nil
}

var illegalFileNameCharacters = regexp.MustCompile(`[^[a-zA-Z0-9]-_]`)

func sanitizeFileNamePart(part string) string {
	part = strings.Replace(part, "/", "-", -1)
	part = illegalFileNameCharacters.ReplaceAllString(part, "")
	return part
}
