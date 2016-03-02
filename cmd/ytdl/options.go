package main

type options struct {
	NoProgress  bool
	Output      string
	Info        bool
	SilentMode  bool
	DebugMode   bool
	Append      bool
	Filters     []string
	Range       string
	DownloadURL bool
	JSON        bool
	Offset      string
}
