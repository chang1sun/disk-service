package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"

	"github.com/gabriel-vasile/mimetype"
)

func GetMIMETypeFromReader(r io.Reader) string {
	t, err := mimetype.DetectReader(r)
	if err != nil {
		log.Fatalf("get file mime type failed, err msg: %v", err)
	}
	if t.String() == "" {
		return "unknown"
	}
	return t.String()
}

func GetFileMD5FromReader(r io.Reader) string {
	m := md5.New()
	_, err := io.Copy(m, r)
	if err != nil {
		log.Fatalf("io copy failed, err msg: %v", err)
	}
	return hex.EncodeToString(m.Sum(nil))
}
