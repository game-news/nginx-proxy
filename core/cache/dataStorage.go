package cache

import (
	"fmt"

	"gamenews.niracler.com/collection/core/meta"
	"gamenews.niracler.com/collection/core/util"
	"github.com/mediocregopher/radix.v2/pool"
)

// PVUV数据存储
func DataStorage(storageChannel chan meta.StorageBlock, redisPool *pool.Pool) {
	for block := range storageChannel {
		prefix := block.CounterType + "_"

		// 逐层增加, 加洋葱皮的过程
		// 维度: 天/小时/分钟
		// 层级: 顶级-大分类-小分类-终极页面
		// 存储模型: Redis SortedSet
		setKeys := []string{
			prefix + block.UData.UrlType + "_day_" + util.GetTime(block.UData.Data.TimeLocal, "day"),
			prefix + block.UData.UrlType + "_hour_" + util.GetTime(block.UData.Data.TimeLocal, "hour"),
			//prefix + block.UData.UrlType + "_min_" + util.GetTime(block.UData.Data.TimeLocal, "min"),
		}

		rowId := block.UData.UrlId

		for _, key := range setKeys {
			ret, err := redisPool.Cmd(block.StorageModel, key, 1, rowId).Int()
			if err != nil || ret <= 0 {
				fmt.Println("DataStorage redis storage error.", block.StorageModel, key, rowId)
				util.Log.Errorln("DataStorage redis storage error.", block.StorageModel, key, rowId)
			}
		}
	}
}

// Click存储
func ClickStorage(storageChannel chan meta.StorageBlock, redisPool *pool.Pool) {
	for block := range storageChannel {
		prefix := block.CounterType + "_"

		key := prefix + block.UData.UrlType
		rowId := block.UData.UrlId

		ret, err := redisPool.Cmd(block.StorageModel, key, 1, rowId).Int()
		if err != nil || ret <= 0 {
			fmt.Println("DataStorage redis storage error.", block.StorageModel, key, rowId)
			util.Log.Errorln("DataStorage redis storage error.", block.StorageModel, key, rowId)
		}
	}
}
