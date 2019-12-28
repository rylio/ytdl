package ytdl

// Itag is a youtube is a static youtube video format
type Itag struct {
	Number        int
	Extension     string
	Resolution    string
	VideoEncoding string
	AudioEncoding string
	AudioBitrate  int
	FPS           int // FPS are frames per second
}

func getItag(itag int) *Itag {
	if itag < len(ITAGS) {
		return ITAGS[itag]
	}
	return nil
}

// ITAGS is a map of all itags and their attributes
var ITAGS = generateItags()

func generateItags() (list []*Itag) {
	list = make([]*Itag, 403)

	add := func(itag Itag) {
		list[itag.Number] = &itag
	}

	add(Itag{
		Number:        5,
		Extension:     "flv",
		Resolution:    "240p",
		VideoEncoding: "Sorenson H.283",
		AudioEncoding: "mp3",
		AudioBitrate:  64,
	})
	add(Itag{
		Number:        6,
		Extension:     "flv",
		Resolution:    "270p",
		VideoEncoding: "Sorenson H.263",
		AudioEncoding: "mp3",
		AudioBitrate:  64,
	})
	add(Itag{
		Number:        13,
		Extension:     "3gp",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
	})
	add(Itag{
		Number:        17,
		Extension:     "3gp",
		Resolution:    "144p",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		AudioBitrate:  24,
	})
	add(Itag{
		Number:        18,
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  96,
	})
	add(Itag{
		Number:        22,
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        34,
		Extension:     "flv",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        35,
		Extension:     "flv",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        36,
		Extension:     "3gp",
		Resolution:    "240p",
		VideoEncoding: "MPEG-4 Visual",
		AudioEncoding: "aac",
		AudioBitrate:  36,
	})
	add(Itag{
		Number:        37,
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        38,
		Extension:     "mp4",
		Resolution:    "3072p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        43,
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        44,
		Extension:     "webm",
		Resolution:    "480p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        45,
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        46,
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        82,
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
		AudioBitrate:  96,
	})
	add(Itag{
		Number:        83,
		Extension:     "mp4",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  96,
	})
	add(Itag{
		Number:        84,
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        85,
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        100,
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        101,
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        102,
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP8",
		AudioEncoding: "vorbis",
		AudioBitrate:  192,
	})

	// DASH (video only)
	add(Itag{
		Number:        133,
		Extension:     "mp4",
		Resolution:    "240p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        134,
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        135,
		Extension:     "mp4",
		Resolution:    "480p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        136,
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        137,
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        138,
		Extension:     "mp4",
		Resolution:    "2160p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        160,
		Extension:     "mp4",
		Resolution:    "144p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        242,
		Extension:     "webm",
		Resolution:    "240p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        243,
		Extension:     "webm",
		Resolution:    "360p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        244,
		Extension:     "webm",
		Resolution:    "480p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        247,
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        248,
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP9",
		AudioBitrate:  9,
	})
	add(Itag{
		Number:        264,
		Extension:     "mp4",
		Resolution:    "1440p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        266,
		Extension:     "mp4",
		Resolution:    "2160p",
		VideoEncoding: "H.264",
	})
	add(Itag{
		Number:        271,
		Extension:     "webm",
		Resolution:    "1440p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        272,
		Extension:     "webm",
		Resolution:    "2160p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        278,
		Extension:     "webm",
		Resolution:    "144p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        298,
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		FPS:           60,
	})
	add(Itag{
		Number:        299,
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		FPS:           60,
	})
	add(Itag{
		Number:        302,
		Extension:     "webm",
		Resolution:    "720p",
		VideoEncoding: "VP9",
		FPS:           60,
	})
	add(Itag{
		Number:        303,
		Extension:     "webm",
		Resolution:    "1080p",
		VideoEncoding: "VP9",
		FPS:           60,
	})
	add(Itag{
		Number:        308,
		Extension:     "webm",
		Resolution:    "1440p",
		VideoEncoding: "VP9",
		FPS:           60,
	})
	add(Itag{
		Number:        313,
		Extension:     "webm",
		Resolution:    "2160p",
		VideoEncoding: "VP9",
	})
	add(Itag{
		Number:        315,
		Extension:     "webm",
		Resolution:    "2160p",
		VideoEncoding: "VP9",
		FPS:           60,
	})

	// DASH (audio only)
	add(Itag{
		Number:        139,
		Extension:     "mp4",
		AudioEncoding: "aac",
		AudioBitrate:  48,
	})
	add(Itag{
		Number:        140,
		Extension:     "mp4",
		AudioEncoding: "aac",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        141,
		Extension:     "mp4",
		AudioEncoding: "aac",
		AudioBitrate:  256,
	})
	add(Itag{
		Number:        171,
		Extension:     "webm",
		AudioEncoding: "vorbis",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        172,
		Extension:     "webm",
		AudioEncoding: "vorbis",
		AudioBitrate:  192,
	})
	add(Itag{
		Number:        249,
		Extension:     "webm",
		AudioEncoding: "opus",
		AudioBitrate:  50,
	})
	add(Itag{
		Number:        250,
		Extension:     "webm",
		AudioEncoding: "opus",
		AudioBitrate:  70,
	})
	add(Itag{
		Number:        251,
		Extension:     "webm",
		AudioEncoding: "opus",
		AudioBitrate:  160,
	})

	// Live streaming
	add(Itag{
		Number:        92,
		Extension:     "ts",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  48,
	})
	add(Itag{
		Number:        93,
		Extension:     "ts",
		Resolution:    "480p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        94,
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        95,
		Extension:     "ts",
		Resolution:    "1080p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  256,
	})
	add(Itag{
		Number:        96,
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  256,
	})
	add(Itag{
		Number:        120,
		Extension:     "flv",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  128,
	})
	add(Itag{
		Number:        127,
		Extension:     "ts",
		AudioEncoding: "aac",
		AudioBitrate:  96,
	})
	add(Itag{
		Number:        128,
		Extension:     "ts",
		AudioEncoding: "aac",
		AudioBitrate:  96,
	})
	add(Itag{
		Number:        132,
		Extension:     "ts",
		Resolution:    "240p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  48,
	})
	add(Itag{
		Number:        151,
		Extension:     "ts",
		Resolution:    "720p",
		VideoEncoding: "H.264",
		AudioEncoding: "aac",
		AudioBitrate:  24,
	})

	add(Itag{
		Number:        394,
		Extension:     "mp4",
		Resolution:    "144p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        395,
		Extension:     "mp4",
		Resolution:    "240p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        396,
		Extension:     "mp4",
		Resolution:    "360p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        397,
		Extension:     "mp4",
		Resolution:    "480p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        398,
		Extension:     "mp4",
		Resolution:    "720p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        399,
		Extension:     "mp4",
		Resolution:    "1080p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        400,
		Extension:     "mp4",
		Resolution:    "1440p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        401,
		Extension:     "mp4",
		Resolution:    "2160p",
		VideoEncoding: "AV1",
	})
	add(Itag{
		Number:        402,
		Extension:     "mp4",
		Resolution:    "2880p",
		VideoEncoding: "AV1",
	})

	return
}
