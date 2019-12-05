package counter

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"nginx-proxy/logcollection/meta"
	"nginx-proxy/logcollection/util"
)

func PvCounter(pvChannel chan meta.UrlData, storageChannel chan meta.StorageBlock) {
	for uData := range pvChannel {
		sItem := meta.StorageBlock{
			CounterType:  "pv",
			StorageModel: "ZINCRBY",
			UData:        uData,
		}
		storageChannel <- sItem
	}
}

func UvCounter(uvChannel chan meta.UrlData, storageChannel chan meta.StorageBlock, redisPool *pool.Pool) {
	for uData := range uvChannel {
		// HyperLogLog redis
		hyperLogLogKey := "uv_hpll_" + util.GetTime(uData.Data.TimeLocal, "day")
		ret, err := redisPool.Cmd("PFADD", hyperLogLogKey, uData.Data.RemoteAddr, "EX", 86400).Int()
		if err != nil {
			util.Log.Warningln("UvCounter check redis hyperloglog failded.", err.Error())
			fmt.Println("UvCounter check redis hyperloglog failded.", err.Error())
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
