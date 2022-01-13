package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func Sha1FromReader(r io.Reader) string {
	_sha1 := sha1.New()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	_sha1.Write(buf.Bytes())
	return hex.EncodeToString(_sha1.Sum(nil))
}
