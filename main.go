package main

import (
	"log"
	"time"

	"gamenews.niracler.com/collection/core/cache"
	"gamenews.niracler.com/collection/core/cache/myredis"
	"gamenews.niracler.com/collection/core/conf"
	"gamenews.niracler.com/collection/core/consumer"
	"gamenews.niracler.com/collection/core/counter"
	"gamenews.niracler.com/collection/core/meta"
)

func main() {
	// 打日志
	// logFd, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY, 0644)
	// if err == nil {
	// 	log.SetOutput(logFd)
	// 	defer logFd.Close()
	// }
	log.Printf("Exec start.\n")
	log.Printf("Params: logFilePath=%s, routineNum=%d, lineNumName=%s", conf.NginxLogFile, conf.RoutineNum, conf.LineNumName)

	// // 初始化一些channel, 用于数据传输
	var logChannel = make(chan string, 3*conf.RoutineNum)
	var pvChannel = make(chan meta.UrlData, conf.RoutineNum)
	var uvChannel = make(chan meta.UrlData, conf.RoutineNum)
	var clickChannel = make(chan meta.UrlData, conf.RoutineNum)
	var pvuvStorChannel = make(chan meta.StorageBlock, conf.RoutineNum)
	// var cliStorChannel = make(chan meta.StorageBlock, conf.RoutineNum)

	// Redis Pool
	redisPool := myredis.RedisPool()

	// 日志消费者
	go consumer.ReadFileLineByLine(logChannel, redisPool)

	// 创建一组日志处理
	for i := 0; i < conf.RoutineNum; i++ {
		go consumer.LogConsumer(logChannel, pvChannel, uvChannel, clickChannel, redisPool)
	}

	// 创建各种统计器
	go counter.PvCounter(pvChannel, pvuvStorChannel)            // pv统计器
	go counter.UvCounter(uvChannel, pvuvStorChannel, redisPool) // uv统计器
	// go counter.ClickCounter(clickChannel, cliStorChannel)       // 点击量统计器
	// TODO: 可以做加更多的统计器

	//创建存储器
	go cache.DataStorage(pvuvStorChannel, redisPool)
	// go cache.ClickStorage(cliStorChannel, redisPool)

	time.Sleep(1000 * time.Hour)
}
