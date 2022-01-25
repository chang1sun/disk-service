package application

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	arepo "github.com/changpro/disk-service/domain/auth/repo"
	aservice "github.com/changpro/disk-service/domain/auth/service"
	"github.com/changpro/disk-service/domain/file/repo"
	frepo "github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/infra/errcode"
	"github.com/changpro/disk-service/infra/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

func QuickUpload(ctx context.Context, userID, fileName, md5 string) (string, error) {
	id, err := tryQuickUpload(ctx, userID, fileName, md5)
	if err != nil {
		return "", err
	}
	if id == "" {
		return "", status.Error(errcode.FindNoFileInServerCode, errcode.FindNoFileInServerMsg)
	}
	return id, nil
}

func tryQuickUpload(ctx context.Context, userID, fileName, md5 string) (string, error) {
	file, err := repo.GetUniFileStoreDao().QueryFileByMd5(ctx, md5)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if file == nil {
		return "", nil
	}
	fileMeta := &repo.UniFileMetaPO{
		Size: file.Length,
		Md5:  md5,
		Type: file.Metadata.Lookup("type").String(),
	}

	// update user's storage size
	err = aservice.UpdateUserStorage(ctx, &arepo.UpdateUserAnalysisDTO{
		UserID:        userID,
		FileNum:       1,
		UploadFileNum: 1,
		Size:          fileMeta.Size,
	})
	if err != nil {
		return "", status.Errorf(errcode.RPCCallErrCode, errcode.RPCCallErrMsg, err)
	}

	// update user's level content
	id, err := repo.GetUserFileDao().AddFileOrDir(ctx,
		buildUploadUserFilePO(file.ID.(primitive.ObjectID).String(), fileName, userID, fileMeta))
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return id, nil
}

// Handle multipart file upload request
func FileUploadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	err := r.ParseMultipartForm(2048)
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
		return
	}
	err = r.ParseForm()
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
		return
	}
	md5 := r.PostForm.Get("md5")
	userID := r.PostForm.Get("user_id")
	f, head, err := r.FormFile("file")
	w.Header().Set("ContentType", "application/json")
	defer f.Close()
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
		return
	}

	// cal md5 and compare
	fileMd5 := util.GetFileMD5FromReader(f)
	log.Printf("cal md5 is: %v\n", fileMd5)
	if fileMd5 != md5 {
		errorResp(errcode.Md5CheckNotPassCode, errcode.Md5CheckNotPassMsg, nil, &w)
		return
	}

	f.Seek(0, io.SeekStart)
	fileName := fmt.Sprintf("%v-%v", head.Filename, time.Now().Format("20060102150405"))
	fileMeta := &frepo.UniFileMetaPO{
		Size:     head.Size,
		Md5:      md5,
		Type:     util.GetMIMETypeFromReader(f),
		UploadBy: userID,
	}
	log.Printf("file type: %v", fileMeta.Type)
	// Write in gridfs
	f.Seek(0, io.SeekStart)
	uId, err := frepo.GetUniFileStoreDao().UploadFile(r.Context(), fileName, f, fileMeta)
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	log.Printf("uid: %v\n", uId)

	// update user's storage size
	err = aservice.UpdateUserStorage(r.Context(), &arepo.UpdateUserAnalysisDTO{
		UserID:        userID,
		FileNum:       1,
		UploadFileNum: 1,
		Size:          head.Size,
	})
	if err != nil {
		s, _ := status.FromError(err)
		errorResp(uint32(s.Code()), s.Message(), nil, &w)
		return
	}

	// update user's level content
	id, err := frepo.GetUserFileDao().AddFileOrDir(r.Context(), buildUploadUserFilePO(uId, fileName, userID, fileMeta))
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}

	w.Write([]byte(fmt.Sprintf(`"file_id: %v", "uni_file_id": %v`, id, uId)))
}

// Handle multipart file upload request
func MPFileUploadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

}

// Handle finish upload request
func FileMergeHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

}

// Handle download request
func DownloadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	err := r.ParseForm()
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
		return
	}
	fileID := r.Form.Get("uniFileId")
	fileName := r.Form.Get("fileName")
	f, err := frepo.GetUniFileStoreDao().GetDownloadStream(r.Context(), fileID)
	defer f.Close()
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	w.Header().Set("Content-Type", f.GetFile().Metadata.Lookup("type").String())
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	log.Println(w.Header())
	// _, err = io.Copy(w, f)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func errorResp(code uint32, msg string, err error, w *http.ResponseWriter) {
	if err == nil {
		(*w).Write([]byte(fmt.Sprintf(`{"code: %v, "msg": %v}`,
			code, msg),
		))
	} else {
		(*w).Write([]byte(fmt.Sprintf(`{"code: %v, "msg": %v}`,
			code, fmt.Sprintf(msg, err)),
		))
	}
}

func buildUploadUserFilePO(uid string, name string, userID string, meta *frepo.UniFileMetaPO) *frepo.UserFilePO {
	return &frepo.UserFilePO{
		UserID:    userID,
		UniFileID: uid,
		Name:      name,
		FileSize:  meta.Size,
		FileMd5:   meta.Md5,
		FileType:  meta.Type,
		IsDir:     2,
		IsHide:    2,
		Path:      "/",
		Status:    1,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}
}
