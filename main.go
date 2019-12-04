package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/sirupsen/logrus"
	"io"
	"nginx-proxy/meta"
	"nginx-proxy/util"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type urlData struct {
	data    meta.DigData
	user    meta.User
	urlType string
	urlId   string
}

type storageBlock struct {
	counterType  string
	storageModel string
	uData        urlData
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

func df(network, addr string) (*redis.Client, error) {
	client, err := redis.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	if err = client.Cmd("AUTH", "123456").Err; err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func main() {
	// 获取参数
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

	// Redis Pool
	redisPool, err := pool.NewCustom("tcp", "music-02.niracler.com:6377", 2*params.routineNum, df)
	if err != nil {
		log.Fatalln("Redis pooll created failed.")
		panic(err)
	} else {
		go func() {
			for {
				redisPool.Cmd("PING")
				time.Sleep(3 * time.Second)
			}
		}()
	}

	// 日志消费者
	go readFileLineByLine(params, logChannel)

	// 创建一组日志处理
	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	// 创建 PV UV 统计器
	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel, redisPool)
	// TODO: 可以做加更多的统计器

	//创建存储器
	go dataStorage(storageChannel, redisPool)
	time.Sleep(1000 * time.Second)
}

// 一行一行地将数据从日志文件中读取到 logChannel 中
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

// 对一行一行的日志进行处理
func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) {
	for logStr := range logChannel {
		// 切割日志字符串, 假如返回的数据是空,那么就不需要解析了
		data := cutLogFetchData(logStr)
		if data == nil {
			continue
		}

		// 获取用户信息
		user := meta.User{}
		claims, err := util.ParseToken(data.HttpToken, []byte("onlinemusic"))
		if err != nil {
			user.IsAnonymous = true
			user.IsAuthenticated = false
			user.IsAdmin = false
		} else {
			user.Uid = int64(claims.(jwt.MapClaims)["uid"].(float64))
			user.Username, _ = claims.(jwt.MapClaims)["name"].(string)
			user.IsAnonymous = false
			user.IsAuthenticated = true
			user.IsAdmin, _ = claims.(jwt.MapClaims)["isAdmin"].(bool)
		}

		// TODO: 可以做更多的处理

		// 将数据放到 Channel
		r1, _ := regexp.Compile("song|author|playlist|user|admin|other|media|static")
		r2, _ := regexp.Compile("/([0-9]+)")
		urlType := r1.FindString(data.HttpUrl)
		if urlType == "" {
			urlType = "other"
		}

		urlId := r2.FindString(data.HttpUrl)
		if urlId != "" {
			urlId = urlId[1:]
		} else {
			urlId = "list"
		}
		uData := urlData{
			data:    *data,
			user:    user,
			urlType: urlType,
			urlId:   urlId,
		}
		pvChannel <- uData
		uvChannel <- uData
	}
}

// 将一行的日志切割到结构体中
func cutLogFetchData(logStr string) *meta.DigData {
	values := strings.Split(logStr, "\"")
	var res []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			res = append(res, value)
		}
	}
	if len(res) > 0 {
		r := strings.Split(res[3], " ")
		if len(r) < 3 {
			log.Warningln("Some different", res[3])
			return nil
		}
		data := meta.DigData{
			RemoteAddr:        res[0],
			RemoteUser:        res[1],
			TimeLocal:         res[2],
			HttpMethod:        r[0],
			HttpUrl:           r[1],
			HttpVersion:       r[2],
			Status:            res[4],
			BodyBytesSent:     res[5],
			HttpReferer:       res[6],
			HttpUserAgent:     res[7],
			HttpXForwardedFor: res[8],
			HttpToken:         res[9],
		}
		return &data
	}

	return nil
}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {
	for uData := range pvChannel {
		sItem := storageBlock{
			counterType:  "pv",
			storageModel: "ZINCRBY",
			uData:        uData,
		}
		storageChannel <- sItem
	}
}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock, redisPool *pool.Pool) {
	for uData := range uvChannel {
		// HyperLogLog redis
		hyperLogLogKey := "uv_hpll_" + getTime(uData.data.TimeLocal, "day")
		ret, err := redisPool.Cmd("PFADD", hyperLogLogKey, uData.data.RemoteAddr, "EX", 86400).Int()
		if err != nil {
			log.Warningln("UvCounter check redis hyperloglog failded.", err.Error())
			fmt.Println("UvCounter check redis hyperloglog failded.", err.Error())
			continue
		}
		if ret != 1 {
			continue
		}

		sItem := storageBlock{
			counterType:  "uv",
			storageModel: "ZINCRBY",
			uData:        uData,
		}
		storageChannel <- sItem
	}
}

// 将日志文件中的时间格式化为时间戳的函数
func getTime(logTime, timeType string) string {
	var item string

	switch timeType {
	case "day":
		item = "2006-01-02"
		break
	case "hour":
		item = "2006-01-02 15"
		break
	case "min":
		item = "2006-01-02 15:04"
		break
	}
	theTime, _ := time.Parse("02/Jan/2006:15:04:05 -0700", logTime)
	t, _ := time.Parse(item, theTime.Format(item))
	return strconv.FormatInt(t.Unix(), 10)
}

// 数据存储
func dataStorage(storageChannel chan storageBlock, redisPool *pool.Pool) {
	for block := range storageChannel {
		prefix := block.counterType + "_"

		// 逐层增加, 加洋葱皮的过程
		// 维度: 天/小时/分钟
		// 层级: 顶级-大分类-小分类-终极页面
		// 存储模型: Redis SortedSet
		setKeys := []string{
			prefix + "day_" + getTime(block.uData.data.TimeLocal, "day"),
			prefix + "hour_" + getTime(block.uData.data.TimeLocal, "hour"),
			prefix + "min_" + getTime(block.uData.data.TimeLocal, "min"),
			prefix + block.uData.urlType + "_day_" + getTime(block.uData.data.TimeLocal, "day"),
			prefix + block.uData.urlType + "_hour_" + getTime(block.uData.data.TimeLocal, "hour"),
			prefix + block.uData.urlType + "_min_" + getTime(block.uData.data.TimeLocal, "min"),
		}

		rowId := block.uData.urlId

		for _, key := range setKeys {
			ret, err := redisPool.Cmd(block.storageModel, key, 1, rowId).Int()
			if err != nil || ret <= 0 {
				fmt.Println("DataStorage redis storage error.", block.storageModel, key, rowId)
				log.Errorln("DataStorage redis storage error.", block.storageModel, key, rowId)
			}
		}
	}
}
