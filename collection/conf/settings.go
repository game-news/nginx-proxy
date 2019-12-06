package conf

import "os"

var LogFilePath = "log/http-access.log"
var LogFile = "log/app.log"
var LineNumName = "log_line_test"
var RedisHost = "music-01.niracler.com:6377"
var RedisPass = "123456"
var ResourceType = "song|author|playlist|user|admin|other|media|static"

func init() {
	// 假如为空，配置默认值
	if os.Getenv("LINE_NUMBER_NAME") != "" {
		LineNumName = os.Getenv("LINE_NUMBER_NAME")
		RedisHost = os.Getenv("REDIS_HOST")
	}
}
