package application

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	arepo "github.com/changpro/disk-service/domain/auth/repo"
	aservice "github.com/changpro/disk-service/domain/auth/service"
	"github.com/changpro/disk-service/domain/file/repo"
	frepo "github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/infra/errcode"
	"github.com/changpro/disk-service/infra/util"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

var MPRedisClient *redis.Client

func tryQuickUpload(ctx context.Context, userID, fileName, md5 string) (string, string, error) {
	file, err := repo.GetUniFileStoreDao().QueryFileByMd5(ctx, md5)
	if err != nil {
		return "", "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if file == nil {
		return "", "", nil
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
		return "", "", status.Errorf(errcode.RPCCallErrCode, errcode.RPCCallErrMsg, err)
	}

	// update user's level content
	id, err := repo.GetUserFileDao().AddFileOrDir(ctx,
		buildUploadUserFilePO(file.ID.(primitive.ObjectID).String(), fileName, userID, fileMeta))
	if err != nil {
		return "", "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return id, file.ID.(primitive.ObjectID).String(), nil
}

// Handle single-piece file upload request
func FileUploadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Header().Set("ContentType", "application/json")
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
	defer f.Close()
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
		return
	}
	// try quick upload
	fid, uniFileID, err := tryQuickUpload(r.Context(), userID, head.Filename, md5)
	if err != nil {
		s, _ := status.FromError(err)
		errorResp(uint32(s.Code()), s.Message(), nil, &w)
		return
	}
	if fid != "" {
		w.Write([]byte(fmt.Sprintf(`"file_id: %v", "uni_file_id": %v`, fid, uniFileID)))
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

// Handle test request by testing if file already exist and acknowledging infos about upcoming formal upload requet
func MPFileUploadTestHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Header().Set("ContentType", "application/json")
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
	fileName := r.PostForm.Get("file_name")
	chunkSize := r.PostForm.Get("chunk_size")
	chunkNum := r.PostForm.Get("chunk_num")
	log.Println("md5 is ", md5)
	// try quick upload
	fid, uniFileID, err := tryQuickUpload(r.Context(), userID, fileName, md5)
	if err != nil {
		s, _ := status.FromError(err)
		errorResp(uint32(s.Code()), s.Message(), nil, &w)
		return
	}
	if fid != "" {
		w.Write([]byte(fmt.Sprintf(`"file_id: %v", "uni_file_id": %v`, fid, uniFileID)))
		return
	}
	// preparation
	// check if file had been uploaded a part of it
	existRes := MPRedisClient.Exists(r.Context(), md5)
	if existRes.Err() != nil {
		log.Fatalf("redis err, err msg: %v", existRes.Err())
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	if existRes.Val() == 1 {
		res := MPRedisClient.HGet(r.Context(), md5, "next_idx")
		if err != nil {
			log.Fatalf("redis err, err msg: %v", res.Err())
			errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
			return
		}
		w.Write([]byte(fmt.Sprintf(`"next_idx: %v"`, res.Val())))
	}
	// initiate
	if err := MPRedisClient.HMSet(r.Context(), md5, "file_name", fileName, "chunk_num",
		chunkNum, "chunk_size", chunkSize, "next_idx", 0).Err(); err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	w.Write([]byte(fmt.Sprintf(`"next_idx: %v"`, 0)))
}

// Handle multipart file upload request
func MPFileUploadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Header().Set("ContentType", "application/json")
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
	chunkID := r.PostForm.Get("chunk_id")
	f, _, err := r.FormFile("file")
	defer f.Close()
	if err != nil {
		errorResp(errcode.ParseHTTPRequestFormFileErrCode, errcode.ParseHTTPRequestFormFileErrMsg, err, &w)
		return
	}
	res := MPRedisClient.HGetAll(r.Context(), md5)
	if res.Err() != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	data := res.Val()
	nextIdx := data["next_idx"]
	chunkNum, _ := strconv.Atoi(data["chunk_num"])
	if chunkID != nextIdx {
		w.Write([]byte(fmt.Sprintf(`"next_idx: %v"`, nextIdx)))
		return
	}
	// if it is the first chunk, set its type
	if nextIdx == "0" {
		MPRedisClient.HSet(r.Context(), md5, "file_type", util.GetMIMETypeFromReader(f))
		f.Seek(0, io.SeekStart)
	}
	// write in temperary file
	tmpDir := fmt.Sprintf("~/tmpfiles/%v", md5)
	if !util.IsPathExists(tmpDir) {
		os.Mkdir(md5, os.ModePerm)
	}
	file, err := os.Create(fmt.Sprintf("~/tmpfiles/%v/chunk_%v", md5, nextIdx))
	if err != nil {
		log.Fatalf("os create tmp file failed, err msg: %v", err)
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	_, err = io.Copy(file, f)
	if err != nil {
		log.Fatalf("os copy file failed, err msg: %v", err)
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	// update next chunk id
	updateIdx, _ := strconv.Atoi(nextIdx)
	updateIdx += 1
	// if it is the last chunk, call merge handler
	if updateIdx == chunkNum {
		MergeFile(r.Context(), w, md5, userID)
		return
	}
	if err := MPRedisClient.HSet(r.Context(), md5, "next_idx", updateIdx).Err(); err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	w.Write([]byte(fmt.Sprintf(`"next_idx: %v"`, updateIdx)))
}

// // Handle finish upload request
// func FileMergeHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

// }

func MergeFile(ctx context.Context, w http.ResponseWriter, md5 string, userID string) {
	res := MPRedisClient.HGetAll(ctx, md5)
	if res.Err() != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, res.Err(), &w)
		log.Fatalf("redis err, err msg: %v", res.Err())
	}
	data := res.Val()
	fileName := data["file_name"]
	chunkNum := data["chunk_num"]
	totalNum, _ := strconv.Atoi(chunkNum)
	file, _ := os.Create(fmt.Sprintf("~/tmpfiles/%v/%v", md5, fileName))
	var pos int64 = 0
	for i := 0; i < totalNum; i++ {
		chunkStream, err := os.OpenFile(fmt.Sprintf("~/tmpfiles/%v/chunk_%v", md5, i), os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatalf("os open file failed, err msg: %v", err)
			errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
			return
		}
		chunkInfo, err := os.Stat(fmt.Sprintf("~/tmpfiles/%v/chunk_%v", md5, i))
		if err != nil {
			log.Fatalf("os open file failed, err msg: %v", err)
			errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
			return
		}
		file.Seek(pos+chunkInfo.Size(), 0)
		_, err = io.Copy(file, chunkStream)
		if err != nil {
			log.Fatalf("os copy file failed, err msg: %v", err)
			errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
			return
		}
	}
	fileType := data["file_type"]
	fileInfo, err := os.Stat(fmt.Sprintf("~/tmpfiles/%v/%v", md5, fileName))
	if err != nil {
		log.Fatalf("os open file failed, err msg: %v", err)
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	fileName = fmt.Sprintf("%v-%v", fileName, time.Now().Format("20060102150405"))
	fileMeta := &frepo.UniFileMetaPO{
		Size:     fileInfo.Size(),
		Md5:      md5,
		Type:     fileType,
		UploadBy: userID,
	}
	log.Printf("file type: %v", fileMeta.Type)
	// Write in gridfs
	file.Seek(0, io.SeekStart)
	uId, err := frepo.GetUniFileStoreDao().UploadFile(ctx, fileName, file, fileMeta)
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	log.Printf("uid: %v\n", uId)

	// update user's storage size
	err = aservice.UpdateUserStorage(ctx, &arepo.UpdateUserAnalysisDTO{
		UserID:        userID,
		FileNum:       1,
		UploadFileNum: 1,
		Size:          fileMeta.Size,
	})
	if err != nil {
		s, _ := status.FromError(err)
		errorResp(uint32(s.Code()), s.Message(), nil, &w)
		return
	}
	// update user's level content
	id, err := frepo.GetUserFileDao().AddFileOrDir(ctx, buildUploadUserFilePO(uId, fileName, userID, fileMeta))
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	w.Write([]byte(fmt.Sprintf(`"file_id: %v", "uni_file_id": %v`, id, uId)))
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
