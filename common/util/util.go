package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"

	"github.com/fatih/structs"
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
	log.Println(s.Map())
	pb, err := structpb.NewStruct(s.Map())
	if err != nil {
		log.Fatalf("cannot convert, err: %v", err)
	}
	return pb
}
