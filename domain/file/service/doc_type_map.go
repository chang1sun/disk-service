package service

var docTypeMap = map[int32][]string{
	1: []string{
		"image/avif",
		"image/bmp",
		"image/gif",
		"image/jpeg",
		"image/png",
		"image/svg+xml",
		"image/webp",
	},
	2: []string{
		"video/x-msvideo",
		"video/mp4",
		"video/mpeg",
		"video/ogg",
	},
	3: []string{
		"audio/aac",
		"application/x-cdf",
		"audio/midi",
		"audio/x-midi",
		"audio/mpeg",
		"audio/ogg",
		"audio/wav",
		"audio/webm",
	},
	4: []string{
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/html",
		"application/pdf",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	},
}
