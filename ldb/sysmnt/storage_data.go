/**
 * 日志仓信息
 */
package sysmnt

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/com"
	"github.com/reggiepy/LogBeetle/global"
	"math"
	"os"
	"time"

	"github.com/shirou/gopsutil/disk"
)

type StorageResult struct {
	Info string          `json:"info,omitempty"`
	Data []*StorageModel `json:"data,omitempty"` // 占用空间
}
type StorageModel struct {
	Name       string `json:"name,omitempty"`       // 名称
	LogCount   uint32 `json:"logCount,omitempty"`   // 日志量
	IndexCount uint32 `json:"indexCount,omitempty"` // 已建索引数量
	FileCount  uint32 `json:"fileCount,omitempty"`  // 文件数量
	TotalSize  string `json:"totalSize,omitempty"`  // 占用空间
}

func init() {
	go func() {
		if global.LbConfig.Store.AutoAddDate && global.LbConfig.Store.SaveDays > 0 {
			// removeStorageByDays() // 注释掉，没必要启动时就清理
			ticker := time.NewTicker(time.Hour) // 一小时检查一次是否有待删除的日志仓
			for {
				<-ticker.C
				removeStorageByDays()
			}
		}
	}()
}

func GetStorageList() *StorageResult {

	var datas []*StorageModel
	var sum int64
	var total int64
	names := com.GetStorageNames(global.LbConfig.Store.Root, ".sysmnt")
	for _, name := range names {
		d := &StorageModel{
			Name: name,
		}

		cnt, size, _ := com.GetDirInfo(global.LbConfig.Store.Root + com.PathSeparator() + name)
		d.TotalSize = com.GetSizeInfo(uint64(size))
		d.FileCount = cnt

		sum += size

		if cnt == 0 {
			d.LogCount = 0
			d.IndexCount = 0
		} else {
			sysmntStore := NewSysmntStorage()
			d.LogCount = sysmntStore.GetStorageDataCount(name)
			d.IndexCount = sysmntStore.GetStorageIndexCount(name)
			total += int64(d.LogCount)
		}

		datas = append(datas, d)
	}

	stat, _ := disk.Usage(global.LbConfig.Store.Root)

	rs := &StorageResult{
		Info: fmt.Sprintf("日志总量 %s 条，共占用空间 %s，剩余空间 %s", com.Int64ToString(total), com.GetSizeInfo(uint64(sum)), com.GetSizeInfo(stat.Free)),
		Data: datas,
	}
	return rs
}

// 删除指定日志仓目录
func DeleteStorage(name string) error {
	// 先尝试目录改名，改成功后再删除
	oldPath := global.LbConfig.Store.Root + com.PathSeparator() + name
	newPath := global.LbConfig.Store.Root + com.PathSeparator() + "_x_" + name
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(newPath)
	if err != nil {
		return err
	}
	return NewSysmntStorage().DeleteStorageInfo(name)
}

func removeStorageByDays() {
	// 日志按日期分仓存储时，按保存天数自动删除
	minYmd := com.GetYyyymmdd(-1 * global.LbConfig.Store.SaveDays)
	dirs := com.GetStorageNames(global.LbConfig.Store.Root, ".sysmnt")
	for _, dir := range dirs {
		ymd := com.Right(dir, 8)
		if com.StringToUint32(ymd, math.MaxUint32) == math.MaxUint32 {
			continue // 后8位不是数字的忽略
		}

		if ymd < minYmd {
			err := DeleteStorage(dir)
			if err != nil {
				global.LbLogger.Info(fmt.Sprintf("日志仓最多保存 %d 天，删除 %s 失败: %v", global.LbConfig.Store.SaveDays, dir, err))
			} else {
				global.LbLogger.Info(fmt.Sprintf("日志仓最多保存 %d 天，已删除 %s", global.LbConfig.Store.SaveDays, dir))
			}
		}
	}
}
