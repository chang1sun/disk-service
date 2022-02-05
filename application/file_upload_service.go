package application

import (
	"context"
	"encoding/json"
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
		return "", "", err
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
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in FileUploadHandler", r)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
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
		log.Fatalf("try quick upload err, err msg: %v", err)
		s, _ := status.FromError(err)
		errorResp(uint32(s.Code()), s.Message(), nil, &w)
		return
	}
	if fid != "" {
		rsp, _ := json.Marshal(map[string]string{"file_id": fid, "uni_file_id": uniFileID})
		w.Write(rsp)
		return
	}
	// cal md5 and compare
	fileMd5 := util.GetFileMD5FromReader(f)
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
	// Write in gridfs
	f.Seek(0, io.SeekStart)
	uId, err := frepo.GetUniFileStoreDao().UploadFile(r.Context(), fileName, f, fileMeta)
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}

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
	id, err := frepo.GetUserFileDao().AddFileOrDir(r.Context(), buildUploadUserFilePO(uId, head.Filename, userID, fileMeta))
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	rsp, _ := json.Marshal(map[string]string{"file_id": id, "uni_file_id": uId})
	w.Write(rsp)
}

// Handle test request by testing if file already exist and acknowledging infos about upcoming formal upload requet
func MPFileUploadTestHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in MPFileUploadTestHandler", r)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
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
	// try quick upload
	fid, uniFileID, err := tryQuickUpload(r.Context(), userID, fileName, md5)
	if err != nil {
		s, _ := status.FromError(err)
		errorResp(uint32(s.Code()), s.Message(), nil, &w)
		return
	}
	if fid != "" {
		rsp, _ := json.Marshal(map[string]string{"file_id": fid, "uni_file_id": uniFileID})
		w.Write(rsp)
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
		nextIdx := res.Val()
		// if all chunks have been uploaded, call merge
		if nextIdx == chunkNum {
			mergeFile(r.Context(), w, md5, userID)
			return
		}
		rsp, _ := json.Marshal(map[string]string{"next_idx": nextIdx})
		log.Println(rsp)
		w.Write(rsp)
		return
	}
	// initiate
	if err := MPRedisClient.HMSet(r.Context(), md5, "file_name", fileName, "chunk_num",
		chunkNum, "chunk_size", chunkSize, "next_idx", 0).Err(); err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	rsp, _ := json.Marshal(map[string]string{"next_idx": "0"})
	w.Write(rsp)
}

// Handle multipart file upload request
func MPFileUploadHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in MPFileUploadHandler", r)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
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
		rsp, _ := json.Marshal(map[string]string{"next_idx": nextIdx})
		log.Println("chunkIdx != nextIdx, rsp is ", rsp)
		w.Write(rsp)
		return
	}
	// if it is the first chunk, set its type
	if nextIdx == "0" {
		MPRedisClient.HSet(r.Context(), md5, "file_type", util.GetMIMETypeFromReader(f))
		f.Seek(0, io.SeekStart)
	}
	// write in temperary file
	tmpDir := fmt.Sprintf("tmpfiles/%v", md5)
	if err = os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	file, err := os.Create(fmt.Sprintf("tmpfiles/%v/chunk_%v", md5, nextIdx))
	if err != nil {
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	_, err = io.Copy(file, f)
	if err != nil {
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	// update next chunk id
	updateIdx, _ := strconv.Atoi(nextIdx)
	updateIdx += 1
	if err := MPRedisClient.HSet(r.Context(), md5, "next_idx", updateIdx).Err(); err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
	// if it is the last chunk, call merge handler
	if updateIdx == chunkNum {
		mergeFile(r.Context(), w, md5, userID)
		return
	}
	rsp, _ := json.Marshal(map[string]string{"next_idx": strconv.Itoa(updateIdx)})
	w.Write(rsp)
}

// // // Handle finish upload request
// func FileMergeHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

// }

func mergeFile(ctx context.Context, w http.ResponseWriter, md5 string, userID string) {
	w.Header().Set("Content-Type", "application/json")
	res := MPRedisClient.HGetAll(ctx, md5)
	if res.Err() != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, res.Err(), &w)
	}
	data := res.Val()
	fileName := data["file_name"]
	chunkNum := data["chunk_num"]
	totalNum, _ := strconv.Atoi(chunkNum)
	file, _ := os.Create(fmt.Sprintf("tmpfiles/%v/%v", md5, fileName))
	var pos int64 = 0
	for i := 0; i < totalNum; i++ {
		chunkStream, err := os.OpenFile(fmt.Sprintf("tmpfiles/%v/chunk_%v", md5, i), os.O_RDONLY, os.ModePerm)
		if err != nil {
			errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
			return
		}
		chunkInfo, err := os.Stat(fmt.Sprintf("tmpfiles/%v/chunk_%v", md5, i))
		if err != nil {
			errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
			return
		}
		file.Seek(pos, 0)
		pos += chunkInfo.Size()
		_, err = io.Copy(file, chunkStream)
		if err != nil {
			log.Fatalf("os copy file failed, err msg: %v", err)
			errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
			return
		}
	}
	fileType := data["file_type"]
	fileInfo, err := os.Stat(fmt.Sprintf("tmpfiles/%v/%v", md5, fileName))
	if err != nil {
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	fileMeta := &frepo.UniFileMetaPO{
		Size:     fileInfo.Size(),
		Md5:      md5,
		Type:     fileType,
		UploadBy: userID,
	}
	// Write in gridfs
	file.Seek(0, io.SeekStart)
	uId, err := frepo.GetUniFileStoreDao().UploadFile(ctx,
		fmt.Sprintf("%v-%v", fileName, time.Now().Format("20060102150405")), file, fileMeta)
	if err != nil {
		errorResp(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err, &w)
		return
	}
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
	// delete tmp dir
	if err := os.RemoveAll(fmt.Sprintf("tmpfiles/%v", md5)); err != nil {
		errorResp(errcode.OsOperationErrCode, errcode.OsOperationErrMsg, err, &w)
		return
	}
	rsp, _ := json.Marshal(map[string]string{"file_id": id, "uni_file_id": uId})
	w.Write(rsp)
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
	log.Fatalf("err occured, code: %v, msg: %v", code, fmt.Sprintf(msg, err))
	if err == nil {
		rsp, _ := json.Marshal(map[string]string{"code": strconv.Itoa(int(code)), "msg": msg})
		(*w).Write(rsp)
	} else {
		rsp, _ := json.Marshal(map[string]string{"code": strconv.Itoa(int(code)), "msg": fmt.Sprintf(msg, err)})
		(*w).Write(rsp)
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
