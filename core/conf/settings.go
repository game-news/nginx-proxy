package conf

import (
	"log"

	"github.com/go-ini/ini"
)

var (
	Cfg          *ini.File
	NginxLogFile string
	LogFile      string
	RoutineNum   int
	LineNumName  string
	RedisHost    string
	RedisPass    string
	ResourceType string
	JWTSecret    string
)

func init() {
	var err error
	path := "core/conf/app.ini"
	Cfg, err = ini.Load(path)
	if err != nil {
		log.Fatalf("Fail to parse '%v': %v", path, err)
	}

	LoadApp()
	LoadDatabase()
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	NginxLogFile = sec.Key("NGINX_LOG_FILE").MustString("log/http-access.log")
	LogFile = sec.Key("LOG_FILE").MustString("log/app.log")
	RoutineNum = sec.Key("ROUNTINE_NUM").MustInt(5)
	LineNumName = sec.Key("LINE_NUMBER_NAME").MustString("log_line")
	ResourceType = sec.Key("RESOURCE_TYPE").MustString("")
	JWTSecret = sec.Key("JWT_SECRET").MustString("")
}

func LoadDatabase() {
	sec, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'database': %v", err)
	}
	RedisHost = sec.Key("REDIS_HOST").MustString("127.0.0.1:6379")
	RedisPass = sec.Key("REDIS_PASS").MustString("123456")
}
