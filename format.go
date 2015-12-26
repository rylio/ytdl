package ytdl

// FormatKey is a string type containg a key in a video format map
type FormatKey string

// Available format Keys
const (
	FormatExtensionKey     FormatKey = "ext"
	FormatResolutionKey    FormatKey = "res"
	FormatVideoEncodingKey FormatKey = "videnc"
	FormatAudioEncodingKey FormatKey = "audenc"
	FormatItagKey          FormatKey = "itag"
)

type Format struct {
	Itag          int    `json:"itag"`
	Extension     string `json: "extension"`
	Resolution    string `json: "resolution"`
	VideoEncoding string `json: "videoEncoding"`
	AudioEncoding string `json: "audioEncoding"`
	meta          map[string]interface{}
}

func newFormat(itag int) (Format, bool) {
	if f, ok := FORMATS[itag]; ok {
		f.meta = make(map[string]interface{})
		return f, true
	}
	return Format{}, false
}

func (f Format) ValueForKey(key FormatKey) interface{} {
	switch key {
	case FormatItagKey:
		return f.Itag
	case FormatExtensionKey:
		return f.Extension
	case FormatResolutionKey:
		return f.Resolution
	case FormatVideoEncodingKey:
		return f.VideoEncoding
	case FormatAudioEncodingKey:
		return f.AudioEncoding
	default:
		return nil
	}
}

// FORMATS is a map of all itags and their formats
var FORMATS = map[int]Format{
	5: Format{
		Extension:     "flv",
		Resolution:    "240p",
		VideoEncoding: "Sorenson H.283",
		AudioEncoding: "mp3",
		Itag:          5,
	},
	6: Format{
		Extension:     "flv",
		Resolution:    "270p",
		VideoEncoding: "Sorenson H.263",
		AudioEncoding: "mp3",
		Itag:          6,
	},
	13: Format{
		Extension:     "3gp",
		Resolution:    "",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		Itag:          13,
	},
	17: Format{
		Extension:     "3gp",
		Resolution:    "144p",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		Itag:          17,
	},
	18: Format{
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          18,
	},
	22: Format{
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          22,
	},
	34: Format{
		Extension:     "flv",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          34,
	},
	35: Format{
		Extension:     "flv",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          35,
	},
	36: Format{
		Extension:     "3gp",
		Resolution:    "240p",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		Itag:          36,
	},
	37: Format{
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          37,
	},
	38: Format{
		Extension:     "mp4",
		Resolution:    "3072p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          38,
	},
	43: Format{
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          43,
	},
	44: Format{
		Extension:     "webm",
		Resolution:    "480p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          44,
	},
	45: Format{
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          45,
	},
	46: Format{
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          46,
	},
	82: Format{
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		Itag:          82,
	},
	83: Format{
		Extension:     "mp4",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          83,
	},
	84: Format{
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          84,
	},
	85: Format{
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          85,
	},
	100: Format{
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          100,
	},
	101: Format{
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          101,
	},
	102: Format{
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		Itag:          102,
	},
	// DASH (video only)
	133: Format{
		Extension:     "mp4",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          133,
	},
	134: Format{
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          134,
	},
	135: Format{
		Extension:     "mp4",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          135,
	},
	136: Format{
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          136,
	},
	137: Format{
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          137,
	},
	138: Format{
		Extension:     "mp4",
		Resolution:    "2160p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          138,
	},
	160: Format{
		Extension:     "mp4",
		Resolution:    "144p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          160,
	},
	242: Format{
		Extension:     "webm",
		Resolution:    "240p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          242,
	},
	243: Format{
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          243,
	},
	244: Format{
		Extension:     "webm",
		Resolution:    "480p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          244,
	},
	247: Format{
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          247,
	},
	248: Format{
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          248,
	},
	264: Format{
		Extension:     "mp4",
		Resolution:    "1440p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          264,
	},
	266: Format{
		Extension:     "mp4",
		Resolution:    "2160p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          266,
	},
	271: Format{
		Extension:     "webm",
		Resolution:    "1440p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          271,
	},
	272: Format{
		Extension:     "webm",
		Resolution:    "2160p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          272,
	},
	278: Format{
		Extension:     "webm",
		Resolution:    "144p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          278,
	},
	298: Format{
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          298,
	},
	299: Format{
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "",
		Itag:          299,
	},
	302: Format{
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          302,
	},
	303: Format{
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP9",
		AudioEncoding: "",
		Itag:          303,
	},
	// DASH (audio only)
	139: Format{
		Extension:     "mp4",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          139,
	},
	140: Format{
		Extension:     "mp4",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          140,
	},
	141: Format{
		Extension:     "mp4",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          141,
	},
	171: Format{
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "vorbis",
		Itag:          171,
	},
	172: Format{
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "vorbis",
		Itag:          172,
	},
	249: Format{
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "opus",
		Itag:          249,
	},
	250: Format{
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "opus",
		Itag:          250,
	},
	251: Format{
		Extension:     "webm",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "opus",
		Itag:          251,
	},
	// Live streaming
	92: Format{
		Extension:     "ts",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          92,
	},
	93: Format{
		Extension:     "ts",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          93,
	},
	94: Format{
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          94,
	},
	95: Format{
		Extension:     "ts",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          95,
	},
	96: Format{
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          96,
	},
	120: Format{
		Extension:     "flv",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          120,
	},
	127: Format{
		Extension:     "ts",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          127,
	},
	128: Format{
		Extension:     "ts",
		Resolution:    "",
		VideoEncoding: "",
		AudioEncoding: "aac",
		Itag:          128,
	},
	132: Format{
		Extension:     "ts",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          132,
	},
	151: Format{
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		Itag:          151,
	},
}
