package com

import (
	"bytes"
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func GeyStoreNameByDate(name string) string {
	if name == "" {
		name = "logdata"
	}
	if global.LbConfig.Store.AutoAddDate {
		return fmt.Sprint(name, "-", time.Now().Format("20060102")) // name-yyyymmdd
	}
	return name
}

func JoinBytes(bts ...[]byte) []byte {
	return bytes.Join(bts, []byte(""))
}

// 取日志仓名列表，以“.”开头的默认忽略
func GetStorageNames(path string, excludes ...string) []string {
	fileinf, err := os.ReadDir(path)
	if err != nil {
		global.LbLogger.Info(fmt.Sprintf("读取目录失败: %v", err))
		return []string{}
	}

	mapDir := make(map[string]string)
	for _, v := range fileinf {
		if v.IsDir() && !Startwiths(v.Name(), ".") {
			mapDir[v.Name()] = ""
		}
	}
	for i := 0; i < len(excludes); i++ {
		delete(mapDir, excludes[i])
	}

	var rs []string
	for k := range mapDir {
		rs = append(rs, k)
	}

	// 倒序
	sort.Slice(rs, func(i, j int) bool {
		return rs[i] > rs[j]
	})

	return rs
}

func GetDirInfo(path string) (uint32, int64, error) {
	var count uint32
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		return err
	})
	return count, size, err
}

func Random() uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		v := r.Uint32()
		if v != 0 {
			return v
		}
	}
}

var localIp string

func GetLocalIp() string {
	if localIp != "" {
		return localIp
	}

	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				localIp = ipnet.IP.String()
			}
		}
	}
	return localIp
}
