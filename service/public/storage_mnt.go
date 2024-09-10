package public

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/com"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/ldb/status"
	"github.com/reggiepy/LogBeetle/ldb/sysmnt"
	"time"
)

type StorageMntService struct{}

// 查询日志仓名称列表
func (s *StorageMntService) Names() []string {
	rs := com.GetStorageNames(global.LbConfig.Store.Root, ".sysmnt")
	return rs
}

// 查询日志仓信息列表
func (s *StorageMntService) List() sysmnt.StorageResult {
	rs := sysmnt.GetStorageList()
	return *rs
}

// 删除指定日志仓
func (s *StorageMntService) Delete(storeName string) error {
	name := storeName
	if name == ".sysmnt" {
		return fmt.Errorf("不能删除 .sysmnt")
	} else if global.LbConfig.Store.AutoAddDate {
		if global.LbConfig.Store.SaveDays > 0 {
			ymd := com.Right(name, 8)
			if com.Len(ymd) == 8 && com.Startwiths(ymd, "20") {
				msg := fmt.Sprintf("当前是日志仓自动维护模式，最多保存 %d 天，不支持手动删除", global.LbConfig.Store.SaveDays)
				return fmt.Errorf(msg)
			}
		}
	} else if name == "logdata" {
		return fmt.Errorf("日志仓 " + name + " 正在使用，不能删除")
	}

	if status.IsStorageOpening(name) {
		return fmt.Errorf("日志仓 " + name + " 正在使用，不能删除")
	}

	err := sysmnt.DeleteStorage(name)
	if err != nil {
		return fmt.Errorf("日志仓 " + name + " 正在使用，不能删除")
	}

	cacheTime = time.Now().Add(-1 * time.Hour) // 让检索时不用缓存名，避免查询不存在的日志仓

	return nil
}
