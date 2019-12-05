package main

import (
	"flag"
	"nginx-proxy/logcollection/cache"
	"nginx-proxy/logcollection/cache/myredis"
	"nginx-proxy/logcollection/consumer"
	"nginx-proxy/logcollection/counter"
	"nginx-proxy/logcollection/meta"
	"nginx-proxy/logcollection/util"
	"os"
	"time"
)

func main() {
	// 获取参数
	logFilePath := flag.String("logFilePath", "log/http-access.log", "log file path")
	routineNum := flag.Int("routineNum", 5, "consumer number by goroutine")
	lineNumName := flag.String("lineNumName", "log_line_1", "consumer number by goroutine")
	l := flag.String("l", "log/app.log", "this program runtime log path")
	flag.Parse()

	params := meta.CmdParams{
		LogFilePath: *logFilePath,
		RoutineNum:  *routineNum,
		LineNumName: *lineNumName,
	}

	// 打日志
	logFd, err := os.OpenFile(*l, os.O_CREATE|os.O_WRONLY, 0644)
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
	var storageChannel = make(chan meta.StorageBlock, params.RoutineNum)

	// Redis Pool
	redisPool := myredis.RedisPool()

	// 日志消费者
	go consumer.ReadFileLineByLine(params, logChannel, redisPool)

	// 创建一组日志处理
	for i := 0; i < params.RoutineNum; i++ {
		go consumer.LogConsumer(params, logChannel, pvChannel, uvChannel, redisPool)
	}

	// 创建 PV UV 统计器
	go counter.PvCounter(pvChannel, storageChannel)
	go counter.UvCounter(uvChannel, storageChannel, redisPool)
	// TODO: 可以做加更多的统计器

	//创建存储器
	go cache.DataStorage(storageChannel, redisPool)
	time.Sleep(1000 * time.Second)
}
