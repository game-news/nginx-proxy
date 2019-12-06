package main

import (
	"fmt"
	"nginx-proxy/collection/cache"
	"nginx-proxy/collection/cache/myredis"
	"nginx-proxy/collection/conf"
	"nginx-proxy/collection/consumer"
	"nginx-proxy/collection/counter"
	"nginx-proxy/collection/meta"
	"nginx-proxy/collection/util"
	"os"
	"time"
)

func main() {
	params := meta.CmdParams{
		LogFilePath: conf.LogFilePath,
		RoutineNum:  5,
		LineNumName: conf.LineNumName,
	}
	fmt.Println(params)

	// 打日志
	logFd, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		util.Log.Out = logFd
		defer logFd.Close()
	}
	util.Log.Infof("Exec start.\n")
	util.Log.Infof("Params: logFilePath=%s, routineNum=%d", params.LogFilePath, params.RoutineNum)

	// 初始化一些channel, 用于数据传输
	var logChannel = make(chan string, 3*params.RoutineNum)
	var pvChannel = make(chan meta.UrlData, params.RoutineNum)
	var uvChannel = make(chan meta.UrlData, params.RoutineNum)
	var clickChannel = make(chan meta.UrlData, params.RoutineNum)
	var pvuvStorChannel = make(chan meta.StorageBlock, params.RoutineNum)
	var cliStorChannel = make(chan meta.StorageBlock, params.RoutineNum)

	// Redis Pool
	redisPool := myredis.RedisPool()

	// 日志消费者
	go consumer.ReadFileLineByLine(params, logChannel, redisPool)

	// 创建一组日志处理
	for i := 0; i < params.RoutineNum; i++ {
		go consumer.LogConsumer(params, logChannel, pvChannel, uvChannel, clickChannel, redisPool)
	}

	// 创建各种统计器
	go counter.PvCounter(pvChannel, pvuvStorChannel)            // pv统计器
	go counter.UvCounter(uvChannel, pvuvStorChannel, redisPool) // uv统计器
	go counter.ClickCounter(clickChannel, cliStorChannel)       // 点击量统计器
	// TODO: 可以做加更多的统计器

	//创建存储器
	go cache.DataStorage(pvuvStorChannel, redisPool)
	go cache.ClickStorage(cliStorChannel, redisPool)

	time.Sleep(1000 * time.Hour)
}
