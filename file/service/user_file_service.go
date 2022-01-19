package service

import (
	"context"
	"time"

	"github.com/changpro/disk-service/common/errcode"
	"github.com/changpro/disk-service/file/repo"
	"github.com/changpro/disk-service/file/rpc/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

const (
	doOverwrite  = 1
	notOverwrite = 2
)

const (
	isDir  = 1
	isFile = 2
)

const (
	isEmpty  = 1
	notEmpty = 2
)

const (
	isHide  = 1
	notHide = 2
)

const (
	statusEnable      = 1
	statusBlackList   = 2
	statusRecycleBin  = 3
	statusDeleted     = 4
	statusPlaceHolder = 5
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
	err = auth.GetAuthCaller().UpdateUserStorage(ctx, userID, 1, 0, fileMeta.Size)
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

func GetUserRoot(ctx context.Context, userID string) ([]*repo.UserFilePO, error) {
	content, err := repo.GetUserFileDao().QueryUserRoot(ctx, userID, false)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return content, nil
}

func GetFileDetail(ctx context.Context, userID, fileID string) (*repo.UserFilePO, error) {
	detail, err := repo.GetUserFileDao().QueryDetail(ctx, fileID)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return detail, nil
}

func MakeNewFolder(ctx context.Context, userID, dirName, path string, overwrite int32) (string, error) {
	// check path correctness
	if err := checkPath(ctx, userID, path); err != nil {
		return "", err
	}
	// check repeat
	if err := checkRepeat(ctx, userID, dirName, path); err != nil {
		return "", err
	}
	id, err := repo.GetUserFileDao().MakeNewFolder(ctx, buildNewFolderPO(userID, dirName, path))
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return id, nil
}

func Rename(ctx context.Context, id, newName string, overwrite int32) error {
	po, err := repo.GetUserFileDao().QueryDetail(ctx, id)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if po == nil {
		return status.Error(errcode.FindNoFileInServerCode, errcode.FindNoFileInServerMsg)
	}
	po.Name = newName
	po.UpdateAt = time.Now()
	if overwrite == doOverwrite {
		_, err = repo.GetUserFileDao().ReplaceFileOrDir(ctx, po)
		if err != nil {
			return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
		}
		return nil
	}
	if err := checkRepeat(ctx, po.UserID, newName, po.Path); err != nil {
		return err
	}
	if err := repo.GetUserFileDao().UpdateFileOrDir(ctx, id,
		&repo.UserFilePO{Name: newName, UpdateAt: time.Now()}); err != nil {
		return err
	}
	return nil
}

func CopyToPath(ctx context.Context, ids []string, path string, overwrite int32) error {
	pos, err := repo.GetUserFileDao().QueryDocByIDs(ctx, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if len(pos) != len(ids) {
		return status.Error(errcode.FindCountNotMatchCode, errcode.FindCountNotMatchMsg)
	}

	// if it is a folder, then trigger recursively call
	for _, po := range pos {
		if po.IsDir == isDir {
			// recursively look up sub folder or files and update them first
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return err
			}
			if len(subPOs) > 0 {
				err = copyToPathByPOs(ctx, subPOs, path+po.Name+"/", overwrite)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, po := range pos {
		po.Path = path
		po.CreateAt = time.Now()
		po.UpdateAt = time.Now()
		po.ID = ""
	}
	err = addDocsToPath(ctx, pos, overwrite)
	if err != nil {
		return err
	}
	return nil
}

func copyToPathByPOs(ctx context.Context, pos []*repo.UserFilePO, path string, overwrite int32) error {
	// if it is a folder, then trigger recursively call
	for _, po := range pos {
		if po.IsDir == isDir {
			// recursively look up sub folder or files and update them first
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return err
			}
			if len(subPOs) > 0 {
				err = copyToPathByPOs(ctx, subPOs, path+po.Name+"/", overwrite)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, po := range pos {
		po.Path = path
		po.CreateAt = time.Now()
		po.UpdateAt = time.Now()
		po.ID = ""
	}
	err := addDocsToPath(ctx, pos, overwrite)
	if err != nil {
		return err
	}
	return nil
}

func replaceToPath(ctx context.Context, po *repo.UserFilePO) (string, error) {
	id, err := repo.GetUserFileDao().ReplaceFileOrDir(ctx, po)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return id, nil
}

func addToPath(ctx context.Context, po *repo.UserFilePO) (string, error) {
	id, err := repo.GetUserFileDao().AddFileOrDir(ctx, po)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return id, nil
}

func addDocsToPath(ctx context.Context, pos []*repo.UserFilePO, overwrite int32) error {
	if overwrite == doOverwrite {
		for _, po := range pos {
			_, err := replaceToPath(ctx, po)
			if err != nil {
				return err
			}
		}
	} else {
		// check repeat
		for _, po := range pos {
			if err := checkRepeat(ctx, po.UserID, po.Name, po.Path); err != nil {
				return err
			}
		}
		// add new
		for _, po := range pos {
			_, err := addToPath(ctx, po)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MoveToPath(ctx context.Context, ids []string, path string, overwrite int32) error {
	// check exist and find po
	pos, err := repo.GetUserFileDao().QueryDocByIDs(ctx, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if len(pos) != len(ids) {
		return status.Error(errcode.FindCountNotMatchCode, errcode.FindCountNotMatchMsg)
	}

	// if it is a folder, then trigger recursively call
	for _, po := range pos {
		if po.IsDir == isDir {
			// recursively look up sub folder or files and update them first
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return err
			}
			if len(subPOs) > 0 {
				err = moveToPathByPOs(ctx, subPOs, path+po.Name+"/", overwrite)
				if err != nil {
					return err
				}
			}
		}
	}

	if overwrite == doOverwrite {
		for _, po := range pos {
			po.Path = path
			po.UpdateAt = time.Now()
			_, err = replaceToPath(ctx, po)
			if err != nil {
				return err
			}
		}
		return nil
	}
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids, &repo.UserFilePO{Path: path, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	return nil
}

func moveToPathByPOs(ctx context.Context, pos []*repo.UserFilePO, path string, overwrite int32) error {
	// if it is a folder, then trigger recursively call
	for _, po := range pos {
		if po.IsDir == isDir {
			// recursively look up sub folder or files and update them first
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return err
			}
			if len(subPOs) > 0 {
				err = moveToPathByPOs(ctx, subPOs, path+po.Name+"/", overwrite)
				if err != nil {
					return err
				}
			}
		}
	}

	if overwrite == doOverwrite {
		for _, po := range pos {
			po.Path = path
			po.UpdateAt = time.Now()
			_, err := replaceToPath(ctx, po)
			if err != nil {
				return err
			}
		}
		return nil
	}
	var ids []string
	for _, po := range pos {
		ids = append(ids, po.ID)
	}
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids, &repo.UserFilePO{Path: path, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	return nil
}

func MoveToRecycleBin(ctx context.Context, ids []string) error {
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{Status: statusRecycleBin, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	return nil
}

func SoftDelete(ctx context.Context, ids []string) error {
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{Status: statusDeleted, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	return nil
}

func HardDelete(ctx context.Context, id string) error {
	if err := repo.GetUserFileDao().DeleteFileOrDir(ctx, id); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func checkPath(ctx context.Context, userID, path string) error {
	ok, err := repo.GetUserFileDao().IsPathExist(ctx, userID, path)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if !ok {
		return status.Error(errcode.PathNotExistCode, errcode.PathNotExistMsg)
	}
	return nil
}

func checkRepeat(ctx context.Context, userID, name, path string) error {
	ok, err := repo.GetUserFileDao().IsFileOrDirExist(ctx, userID, name, path)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if ok {
		return status.Error(errcode.DirOrFileAlreadyExistCode, errcode.DirOrFileAlreadyExistMsg)
	}
	return nil
}

func buildNewFolderPO(userID, name, path string) *repo.UserFilePO {
	return &repo.UserFilePO{
		UserID:   userID,
		Name:     name,
		Path:     path,
		IsDir:    isDir,
		IsHide:   notHide,
		Status:   statusEnable,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
}
