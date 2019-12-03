package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type digData struct {
	time  string
	url   string
	refer string
	us    string
}

type urlData struct {
	data     digData
	username string
}

type urlNode struct {
}

type storageBlock struct {
	counterType  string
	storageModel string
	unode        urlNode
}

type cmdParams struct {
	logFilePath string
	routineNum  int
}

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	// 获取参数\
	logFilePath := flag.String("logFilePath", "log/http-access.log", "log file path")
	routineNum := flag.Int("routineNum", 5, "consumer number by goroutine")
	l := flag.String("l", "log/app.log", "this program runtime log path")
	flag.Parse()

	params := cmdParams{
		logFilePath: *logFilePath,
		routineNum:  *routineNum,
	}

	// 打日志
	logFd, err := os.OpenFile(*l, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	}
	log.Infof("Exec start.\n")
	log.Infof("Params: logFilePath=%s, routineNum=%d", params.logFilePath, params.routineNum)

	// 初始化一些channel, 用于数据传输
	var logChannel = make(chan string, 3*params.routineNum)
	var pvChannel = make(chan urlData, params.routineNum)
	var uvChannel = make(chan urlData, params.routineNum)
	var storageChannel = make(chan storageBlock, params.routineNum)

	// 日志消费者
	go readFileLineByLine(params, logChannel)

	// 创建一组日志处理
	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	// 创建 PV UV 统计器
	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel)
	// 可扩展的

	//创建存储器
	go dataStorage(storageChannel)
	time.Sleep(1000 * time.Second)
}

func readFileLineByLine(params cmdParams, logChannel chan string) {
	fd, err := os.Open(params.logFilePath)
	if err != nil {
		log.Warn("ReadFileLineByLine can not open file, err:" + err.Error())
		return
	}
	defer fd.Close()

	count := 0
	bufferRead := bufio.NewReader(fd)
	for {
		line, err := bufferRead.ReadString('\n')
		logChannel <- line
		log.Infof("line: ", line)
		fmt.Printf(line)
		count++

		if count%(1000*params.routineNum) == 0 {
			log.Infof("ReadFileLineByLine line: %d", count)
		}
		if err != nil {
			if err == io.EOF {
				time.Sleep(3 * time.Second)
				log.Infof("ReadFileLineByLine wait, readline: %d", count)
			} else {
				log.Warn("ReadFileLineByLine read log err :" + err.Error())
			}
		}
	}
}

func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) {

}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {

}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock) {

}

func dataStorage(storageChannel chan storageBlock) {

}
