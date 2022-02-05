package service

import (
	"context"
	"time"

	"github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/infra/errcode"
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

const (
	displayHidden    = 1
	notDisplayHidden = 2
)

func GetUserRoot(ctx context.Context, userID string) ([]*repo.UserFilePO, error) {
	content, err := repo.GetUserFileDao().QueryUserRoot(ctx, userID, false)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return content, nil
}

func GetDirByPath(ctx context.Context, userID, path string, showHide bool) ([]*repo.UserFilePO, error) {
	content, err := repo.GetUserFileDao().QueryDirByPath(ctx, userID, path, showHide)
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
	if err := CheckPath(ctx, userID, path); err != nil {
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

func SetHiddenDoc(ctx context.Context, ids []string, hideStatus int32) error {
	if len(ids) == 0 {
		return nil
	}
	// recursive call
	pos, err := repo.GetUserFileDao().QueryDocByIDs(ctx, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if len(pos) != len(ids) {
		return status.Error(errcode.FindCountNotMatchCode, errcode.FindCountNotMatchMsg)
	}
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = setHiddenDocByPOs(ctx, subPOs, hideStatus)
				if err != nil {
					return err
				}
			}
		}
	}

	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{IsHide: hideStatus, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	return nil
}

func setHiddenDocByPOs(ctx context.Context, pos []*repo.UserFilePO, hideStatus int32) error {
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = setHiddenDocByPOs(ctx, subPOs, hideStatus)
				if err != nil {
					return err
				}
			}
		}
	}
	var ids []string
	for _, po := range pos {
		ids = append(ids, po.ID)
	}
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{IsHide: hideStatus, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	return nil
}

func MoveToRecycleBin(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	// recursive call
	pos, err := repo.GetUserFileDao().QueryDocByIDs(ctx, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if len(pos) != len(ids) {
		return status.Error(errcode.FindCountNotMatchCode, errcode.FindCountNotMatchMsg)
	}
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = moveToRecycleBinByPOs(ctx, subPOs)
				if err != nil {
					return err
				}
			}
		}
	}
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{Status: statusRecycleBin, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	// only insert first level doc to recycle-bin
	err = repo.GetRecycleFileDao().InsertToRecycleBin(ctx, userID, buildRecyclePOs(pos))
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func buildRecyclePOs(list []*repo.UserFilePO) []*repo.RecycleFilePO {
	var res []*repo.RecycleFilePO
	for _, po := range list {
		res = append(res, &repo.RecycleFilePO{
			ID:       po.ID,
			UserID:   po.UserID,
			Name:     po.Name,
			IsDir:    po.IsDir,
			DeleteAt: time.Now(),
			ExpireAt: time.Now().Add(7 * 24 * time.Hour),
		})
	}
	return res
}

func moveToRecycleBinByPOs(ctx context.Context, pos []*repo.UserFilePO) error {
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = moveToRecycleBinByPOs(ctx, subPOs)
				if err != nil {
					return err
				}
			}
		}
	}
	var ids []string
	for _, po := range pos {
		ids = append(ids, po.ID)
	}
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

func SoftDelete(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	// recursive call
	pos, err := repo.GetUserFileDao().QueryDocByIDs(ctx, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if len(pos) != len(ids) {
		return status.Error(errcode.FindCountNotMatchCode, errcode.FindCountNotMatchMsg)
	}
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = softDeleteByPOs(ctx, subPOs)
				if err != nil {
					return err
				}
			}
		}
	}
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{Status: statusDeleted, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	// only delete first level doc in recycle-bin
	err = repo.GetRecycleFileDao().DeleteDocs(ctx, userID, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func softDeleteByPOs(ctx context.Context, pos []*repo.UserFilePO) error {
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = softDeleteByPOs(ctx, subPOs)
				if err != nil {
					return err
				}
			}
		}
	}
	var ids []string
	for _, po := range pos {
		ids = append(ids, po.ID)
	}
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
	po, err := repo.GetUserFileDao().QueryDetail(ctx, id)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if po.IsDir == 1 {
		return status.Error(errcode.ForbidHardDeleteFolderCode, errcode.ForbidHardDeleteFolderMsg)
	}
	if err := repo.GetUserFileDao().DeleteFileOrDir(ctx, id); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func CheckPath(ctx context.Context, userID, path string) error {
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

func GetDirSizeAndSubFilesNum(ctx context.Context, dir *repo.UserFilePO) (int64, int32, error) {
	// when meet file, jump out
	if dir.IsDir == isFile {
		return dir.FileSize, 1, nil
	}
	// recursively cal sub folder
	var totalSize int64
	var totalNum int32
	subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, dir.UserID, dir.Path+dir.Name+"/", true)
	if err != nil {
		return 0, 0, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	for _, subPO := range subPOs {
		curSize, curNum, err := GetDirSizeAndSubFilesNum(ctx, subPO)
		if err != nil {
			return 0, 0, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
		}
		totalSize += curSize
		totalNum += curNum

	}
	return totalSize, totalNum, nil
}

func GetRecycleBinList(ctx context.Context, userID string, offset, limit int32) ([]*repo.RecycleFilePO, error) {
	list, err := repo.GetRecycleFileDao().GetRecycleBinList(ctx, userID, offset, limit)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return list, nil
}

func RecoverDocs(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	// recursive call
	pos, err := repo.GetUserFileDao().QueryDocByIDs(ctx, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if len(pos) != len(ids) {
		return status.Error(errcode.FindCountNotMatchCode, errcode.FindCountNotMatchMsg)
	}
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = recoverDocsByPOs(ctx, subPOs)
				if err != nil {
					return err
				}
			}
		}
	}
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{Status: statusRecycleBin, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	// only delete first level doc in recycle-bin
	err = repo.GetRecycleFileDao().DeleteDocs(ctx, userID, ids)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func recoverDocsByPOs(ctx context.Context, pos []*repo.UserFilePO) error {
	for _, po := range pos {
		if po.IsDir == isDir {
			subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
			if err != nil {
				return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
			}
			if len(subPOs) > 0 {
				err = recoverDocsByPOs(ctx, subPOs)
				if err != nil {
					return err
				}
			}
		}
	}
	var ids []string
	for _, po := range pos {
		ids = append(ids, po.ID)
	}
	count, err := repo.GetUserFileDao().UpdateFileOrDirByIDs(ctx, ids,
		&repo.UserFilePO{Status: statusEnable, UpdateAt: time.Now()})
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count != len(ids) {
		return status.Error(errcode.UpdateCountNotMatchCode, errcode.UpdateCountNotMatchMsg)
	}
	return nil
}
