package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/changpro/disk-service/common/errcode"
	cutil "github.com/changpro/disk-service/common/util"
	"github.com/changpro/disk-service/file/repo"
	"github.com/changpro/disk-service/file/util"
)

// Handle multipart file upload request
func FileUploadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Header().Set("ContentType", "application/json")
	err := r.ParseForm()
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
	}
	md5 := r.Form.Get("md5")
	userID := r.Form.Get("user_id")
	f, head, err := r.FormFile("file")
	defer f.Close()
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
	}
	fileName := fmt.Sprintf("%v-%v", head.Filename, time.Now().Format("20060102150405"))
	fileMeta := &repo.UniFileMetaPO{
		Size:     head.Size,
		Md5:      md5,
		Type:     util.GetFileTypeFromReader(f),
		UploadAt: time.Now(),
		UploadBy: userID,
	}

	// cal md5 and compare
	fileMd5 := cutil.Sha1FromReader(f)
	if fileMd5 != fileMeta.Md5 {
		errorResp(errcode.Md5CheckNotPassCode, errcode.Md5CheckNotPassMsg, err, &w)
	}

	// Write in gridfs
	id, err := repo.GetUniFileStoreDao().UploadFile(r.Context(), fileName, f, fileMeta)
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
	}
	w.Write([]byte(fmt.Sprintf(`"file_id": %v`, id)))
}

// Handle multipart file upload request
func MPFileUploadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

}

// Handle finish upload request
func FileMergeHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

}

func errorResp(code uint32, msg string, err error, w *http.ResponseWriter) {
	(*w).Write([]byte(fmt.Sprintf(`{"code: %v, "msg": %v}`,
		code, msg, err),
	))
}
