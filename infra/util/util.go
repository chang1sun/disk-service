package util

import (
	"bytes"
	"crypto/md5"
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
		log.Fatalf("cannot convert, err: %v", err)
	}
	return pb
}

func GetStringWithSalt(s string) string {
	return Sha1([]byte(s + config.GetConfig().PwSalt))
}

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

func GetMd5FromJson(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

// 判断所给路径文件/文件夹是否存在
func IsPathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
