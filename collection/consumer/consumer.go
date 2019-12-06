package consumer

import (
	"bufio"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mediocregopher/radix.v2/pool"
	"io"
	"nginx-proxy/collection/conf"
	"nginx-proxy/collection/meta"
	"nginx-proxy/collection/util"
	"os"
	"regexp"
	"strings"
	"time"
)

// 一行一行地将数据从日志文件中读取到 logChannel 中
func ReadFileLineByLine(params meta.CmdParams, logChannel chan string, redisPool *pool.Pool) {
	maxNum, err := redisPool.Cmd("GET", params.LineNumName).Int()
	if err != nil {
		maxNum = 0
	}

	fd, err := os.Open(params.LogFilePath)
	if err != nil {
		util.Log.Warn("ReadFileLineByLine can not open file, err:" + err.Error())
		return
	}
	defer fd.Close()

	count := 0
	bufferRead := bufio.NewReader(fd)
	for {
		line, err := bufferRead.ReadString('\n')
		count++

		if maxNum > count {
			continue
		}
		logChannel <- line

		if count%(1000*params.RoutineNum) == 0 {
			util.Log.Infof("ReadFileLineByLine line: %d", count)
		}
		if err != nil {
			if err == io.EOF {
				time.Sleep(3 * time.Second)
				util.Log.Infof("ReadFileLineByLine wait, readline: %d", count)
			} else {
				util.Log.Warn("ReadFileLineByLine read log err :" + err.Error())
			}
		}
	}
}

// 对一行一行的日志进行处理
func LogConsumer(params meta.CmdParams, logChannel chan string, pvChannel, uvChannel, clickChannel chan meta.UrlData, redisPool *pool.Pool) {
	for logStr := range logChannel {
		// 切割日志字符串, 假如返回的数据是空,那么就不需要解析了
		data := cutLogFetchData(logStr)
		if data == nil {
			continue
		}

		_, err := redisPool.Cmd("INCR", params.LineNumName).Int()
		if err != nil {
			fmt.Println(err.Error())
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
		if !user.IsAnonymous {
			key := "action_" + fmt.Sprintf("%d", user.Uid)
			_, err := redisPool.Cmd("LPUSH", key, logStr).Int()
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		// 将数据放到 Channel
		r1, _ := regexp.Compile(conf.ResourceType)
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
		uData := meta.UrlData{
			Data:    *data,
			User:    user,
			UrlType: urlType,
			UrlId:   urlId,
		}
		pvChannel <- uData
		uvChannel <- uData
		clickChannel <- uData
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
			util.Log.Warningln("Some different", res[3])
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
