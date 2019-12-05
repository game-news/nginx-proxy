package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	_ = filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		result = info.Size()
		return nil
	})
	return result
}

// 将日志文件中的时间格式化为时间戳的函数
func GetTime(logTime, timeType string) string {
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

var Log = logrus.New()

func init() {
	Log.Out = os.Stdout
	Log.SetLevel(logrus.DebugLevel)
}
