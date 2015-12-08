package ytdl

// From http://stackoverflow.com/questions/2068344/how-do-i-get-a-youtube-video-thumbnail-from-the-youtube-api

// ThumbnailQuality is a youtube thumbnail quality option
type ThumbnailQuality string

// ThumbnailQualityHigh is the high quality thumbnail jpg
const ThumbnailQualityHigh ThumbnailQuality = "hqdefault"

// ThumbnailQualityDefault is the default quality thumbnail jpg
const ThumbnailQualityDefault ThumbnailQuality = "default"

// ThumbnailQualityMedium is the medium quality thumbnail jpg
const ThumbnailQualityMedium ThumbnailQuality = "mqdefault"

// ThumbnailQualitySD is the standard def quality thumbnail jpg
const ThumbnailQualitySD ThumbnailQuality = "sddefault"

// ThumbnailQualityMaxRes is the maximum resolution quality jpg
const ThumbnailQualityMaxRes ThumbnailQuality = "maxresdefault"
