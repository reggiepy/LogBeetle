package public

import (
	"github.com/reggiepy/LogBeetle/com"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/ldb"
	"github.com/reggiepy/LogBeetle/ldb/search"
	"github.com/reggiepy/LogBeetle/ldb/storage/logdata"
	"github.com/reggiepy/LogBeetle/ldb/sysmnt"
	"time"
)

type LogService struct{}

type storageItem struct {
	storeName     string // 日志仓
	total         uint32 // 日志件数
	isSearchRange bool   // 是否条件范围的日志仓
}

var cacheStoreNames []string // 所有的日志仓（避免每次读磁盘，适当使用缓存）
var cacheTime time.Time      // 最近一次读日志仓目录的时间点

// 日志检索（表单提交方式）
func (s *LogService) Search(cond *search.SearchCondition) *search.SearchResult {
	// 准备好各种场景的检索条件（系统【~】、日志级别【!】、用户【@】）
	startTime := time.Now()
	mnt := sysmnt.NewSysmntStorage()
	cond.SearchSize = global.LbConfig.Search.PageSize
	if cond.User != "" {
		cond.User = "@" + cond.User // 有指定用户条件
	}
	if len(cond.Loglevels) <= 1 || len(cond.Loglevels) >= 4 {
		cond.Loglevels = make([]string, 0) // 多选的单选或全选，都清空（单选走loglevel索引，全选等于没选）
	}
	if !com.IsBlank(cond.Loglevel) && !com.Contains(cond.Loglevel, ",") {
		cond.Loglevel = "!" + com.Trim(cond.Loglevel) // 编辑日志级别单选条件，以便精确匹配
	} else {
		cond.Loglevel = "" // 清空日志级别单选条件，以便多选配配（改用loglevels）
	}

	cond.OrgSystems = append(cond.OrgSystems, "*") // 不需登录时全部系统都有访问权限
	if cond.OrgSystem != "" {
		cond.OrgSystem = "~" + cond.OrgSystem // 多个系统权限，按输入的系统作条件
	}

	// 范围内的日志仓都查一遍
	// 注1）日志不断新增时，总件数可能会因为时间点原因不适最新，从而变现出点点小误差【完全可接受】
	// 注2）跨仓检索时，非本次检索的目标仓的话，只查取相关件数不做真正筛选计数以提高性能，最大匹配件数有时可能出现较大误差【折中可接受】
	result := &search.SearchResult{PageSize: com.IntToString(global.LbConfig.Search.PageSize)}
	var total uint32
	var count uint32
	storeItems := getStoreItems(cond.StoreName, cond.DatetimeFrom, cond.DatetimeTo)
	for i, max := 0, len(storeItems); i < max; i++ {
		item := storeItems[i]
		if !item.isSearchRange {
			// 不需要查数据，只查关联件数
			total += mnt.GetStorageDataCount(item.storeName) // 累加总件数
			continue
		}

		cond.SearchSize = global.LbConfig.Search.PageSize - len(result.Data) // 本次需要查多少件
		if cond.CurrentStoreName != "" && item.storeName > cond.CurrentStoreName {
			cond.SearchSize = 0 // 是范围内的日志仓，但不是本次要查的，设为0不查数据，只查关联件数
		}

		eng := ldb.NewEngine(item.storeName)     // 遍历日志仓检索
		rs := eng.Search(cond)                   // 【检索】按动态的要求件数检索
		total += com.StringToUint32(rs.Total, 0) // 累加总件数
		count += com.StringToUint32(rs.Count, 0) // 累加最大匹配件数
		if len(rs.Data) > 0 {
			result.Data = append(result.Data, rs.Data...) // 累加查询结果
			result.LastStoreName = item.storeName         // 设定检索结果最后一条（最久远）日志所在的日志仓，页面向下滚动继续检索时定位用
		}

		if !(cond.CurrentStoreName != "" && item.storeName > cond.CurrentStoreName) {
			// 仅针对更久远的日志仓
			if len(result.Data) < global.LbConfig.Search.PageSize && i < max-1 {
				// 数据没查够，且后面还有日志仓待查询，准备好跨仓查询条件
				cond.CurrentId = 0         // 下一日志仓从头开始查
				cond.CurrentStoreName = "" // 从头开始所以这个条件不再适用，清空
			}
		}
	}

	// 返回结果
	result.Total = com.Uint32ToString(total)                                      // 总件数
	result.Count = com.Uint32ToString(count)                                      // 最大匹配检索（笼统，在最大查取件数（5000件）内查完时，前端会改成精确的和结果一样的件数）
	result.TimeMessage = "耗时" + getTimeInfo(time.Since(startTime).Milliseconds()) // 查询耗时
	if cond.NewNearId > 0 && len(result.Data) == 0 {
		// 相邻检索时确保能返回定位日志
		result.Data = append(result.Data, search.GetLogDataModelById(cond.NearStoreName, cond.NewNearId))
	}
	return result
}

// 筛选出日志仓检索范围
func getStoreItems(storeName string, datetimeFrom string, datetimeTo string) []*storageItem {
	sysmntStore := sysmnt.NewSysmntStorage()
	var items []*storageItem
	if !global.LbConfig.Store.AutoAddDate {
		// 单日志仓
		name := com.GeyStoreNameByDate("")
		items = append(items, &storageItem{storeName: name, total: sysmntStore.GetStorageDataCount(name), isSearchRange: true})
		return items
	}

	// 遍历日志仓，比较日期范围筛选日志仓
	hasDateCond := (datetimeFrom != "" && datetimeTo != "")     // 是否有日期范围条件
	from := com.ReplaceAll(com.Left(datetimeFrom, 10), "-", "") // yyyymmdd或“”
	to := com.ReplaceAll(com.Left(datetimeTo, 10), "-", "")     // yyyymmdd或“”
	if time.Since(cacheTime) >= time.Second*10 {
		cacheStoreNames = com.GetStorageNames(global.LbConfig.Store.Root, ".sysmnt") // 所有的日志仓，结果已排序，缓存10秒避免频繁读盘
		cacheTime = time.Now()
	}
	for i, max := 0, len(cacheStoreNames); i < max; i++ {
		name := cacheStoreNames[i]
		item := &storageItem{storeName: name, total: sysmntStore.GetStorageDataCount(name)}
		date := com.Right(name, 8) // yyyymmdd
		if storeName == "" {
			// 日志仓条件空白
			if hasDateCond {
				if date >= from && date <= to {
					item.isSearchRange = true // 日期范围内的日志仓都是条件范围
				}
			} else {
				item.isSearchRange = true // 无日志仓条件、且无日期条件，全部都是条件范围了
			}
		} else {
			// 有日志仓条件
			if hasDateCond {
				if storeName == name && date >= from && date <= to {
					item.isSearchRange = true // 有日期条件，得满足日期条件，该日志仓才是条件范围
				}
			} else {
				if storeName == name {
					item.isSearchRange = true // 没日期条件，仅该日志仓是条件范围
				}
			}
		}
		items = append(items, item)
	}
	return items
}

func getTimeInfo(milliseconds int64) string {
	if milliseconds >= 1000 {
		return " " + com.Float64ToString(com.Round1(float64(milliseconds)/1000.0)) + " 秒"
	}
	return " " + com.Int64ToString(milliseconds) + " 毫秒"
}

// 添加2000条测试日志（仅测试模式有效），用于测试或快速体验
func (s *LogService) AddTestData() error {

	//if !conf.IsTestMode() {
	//	return fmt.Errorf("当前不是测试模式，不支持生成测试数据") // 非测试模式时忽略
	//}

	cnt := 0
	for {
		cnt++
		traceId := com.RandomHashString()
		md := &logdata.LogDataModel{
			Text:       "测试用的日志，字段名为Text，" + "字段Date的格式为YYYY-MM-DD HH:MM:SS.SSS，必须格式一致才能正常使用时间范围检索条件。" + "随机3位字符串：" + com.RandomString(3) + "，第" + com.IntToString(cnt) + "条",
			Date:       com.Now(),
			System:     "demo1",
			ServerName: "default",
			ServerIp:   "127.0.0.1",
			ClientIp:   "127.0.0.1",
			TraceId:    traceId,
			LogLevel:   "INFO",
			User:       "tuser-" + com.RandomString(1),
		}
		addDataModelLog(md)

		md2 := &logdata.LogDataModel{
			Text:       "几个随机字符串供查询试验：" + com.RandomString(1) + "，" + com.Right(com.ULID(), 2) + "，" + com.RandomString(3) + "，" + com.Right(com.ULID(), 4) + "，" + com.Right(com.ULID(), 5),
			Date:       com.Now(),
			System:     "demo2",
			ServerName: "default",
			ServerIp:   "127.0.0.1",
			ClientIp:   "127.0.0.1",
			TraceId:    traceId,
			LogLevel:   "DEBUG",
			User:       "tuser-" + com.RandomString(1),
		}
		addDataModelLog(md2)

		md3 := &logdata.LogDataModel{
			Text:       "几个随机字符串供查询试验：" + com.RandomString(1) + "，" + com.Right(com.ULID(), 2) + "，" + com.RandomString(3) + "，" + com.Right(com.ULID(), 4) + "，" + com.Right(com.ULID(), 5),
			Date:       com.Now(),
			System:     "demo3",
			ServerName: "default",
			ServerIp:   "127.0.0.1",
			ClientIp:   "127.0.0.1",
			TraceId:    traceId,
			LogLevel:   "WARN",
			User:       "tuser-" + com.RandomString(1),
		}
		addDataModelLog(md3)

		if cnt >= 1000 {
			break
		}
	}
	return nil
}

// 添加日志（JSON提交方式）
func (s *LogService) JsonLogAdd(md *logdata.LogDataModel) error {
	md.Text = com.Trim(md.Text)
	addDataModelLog(md)
	return nil
}

// 添加日志
func addDataModelLog(data *logdata.LogDataModel) {
	engine := ldb.NewDefaultEngine()
	engine.AddLogDataModel(data)
}
