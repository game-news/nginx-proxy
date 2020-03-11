package counter

import (
	"gamenews.niracler.com/collection/core/meta"
	"gamenews.niracler.com/collection/core/util"
	"github.com/mediocregopher/radix.v2/pool"
)

// 资源点击量计算, "click_" + resType
func ClickCounter(clickChannel chan meta.UrlData, storageChannel chan meta.StorageBlock) {
	for uData := range clickChannel {
		sItem := meta.StorageBlock{
			CounterType:  "click",
			StorageModel: "ZINCRBY",
			UData:        uData,
		}
		storageChannel <- sItem
	}
}

// 浏览量计算
func PvCounter(pvChannel chan meta.UrlData, pvuvStorChannel chan meta.StorageBlock) {
	for uData := range pvChannel {
		sItem := meta.StorageBlock{
			CounterType:  "pv",
			StorageModel: "ZINCRBY",
			UData:        uData,
		}
		pvuvStorChannel <- sItem
	}
}

// 用户量计算
func UvCounter(uvChannel chan meta.UrlData, storageChannel chan meta.StorageBlock, redisPool *pool.Pool) {
	for uData := range uvChannel {
		// HyperLogLog redis
		hyperLogLogKey := "hpll_" + util.GetTime(uData.Data.TimeLocal, "day")
		ret, err := redisPool.Cmd("PFADD", hyperLogLogKey, uData.Data.RemoteAddr, "EX", 86400).Int()
		if err != nil {
			util.Log.Warningln("UvCounter check redis hyperloglog failded.", err.Error())
			continue
		}
		if ret != 1 {
			continue
		}

		sItem := meta.StorageBlock{
			CounterType:  "uv",
			StorageModel: "ZINCRBY",
			UData:        uData,
		}
		storageChannel <- sItem
	}
}
