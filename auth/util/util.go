package util

import (
	"github.com/changpro/disk-service/auth/config"
	cutil "github.com/changpro/disk-service/common/util"
)

func GetStringWithSalt(s string) string {
	return cutil.Sha1([]byte(s + config.GetConfig().PwSalt))
}
