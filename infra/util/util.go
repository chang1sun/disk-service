package util

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"

	"github.com/changpro/disk-service/infra/config"
	"github.com/fatih/structs"
	"github.com/gabriel-vasile/mimetype"
	"google.golang.org/protobuf/types/known/structpb"
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

func AnyToStructpb(i interface{}) *structpb.Struct {
	s := structs.New(i)
	s.TagName = "json"
	pb, err := structpb.NewStruct(s.Map())
	if err != nil {
		log.Printf("cannot convert, err: %v", err)
	}
	return pb
}

func GetStringWithSalt(s string) string {
	return Sha1([]byte(s + config.GetConfig().PwSalt))
}

func GetMIMETypeFromReader(r io.Reader) string {
	t, err := mimetype.DetectReader(r)
	if err != nil {
		log.Printf("get file mime type failed, err msg: %v", err)
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
		log.Printf("io copy failed, err msg: %v", err)
	}
	return hex.EncodeToString(m.Sum(nil))
}

func GetMd5FromJson(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

const (
	StdLen  = 16
	UUIDLen = 20
)

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// var AsciiChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

func NewRandomString() string {
	return NewLenChars(StdLen, StdChars)
}

func NewLenRandomString(length int) string {
	return NewLenChars(length, StdChars)
}

func NewLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
