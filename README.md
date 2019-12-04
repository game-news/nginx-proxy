# nginx-proxy

用户动态

## go 安装包时需要配置的代理

```shell script
git config --global http.proxy socks5://127.0.0.1:1080
export http_proxy=socks5://127.0.0.1:1080
```

## 如何团队项目保持同步(重要)

([附上IDEA可视化操作](https://blog.csdn.net/autfish/article/details/52513465))

第一次时需要,与团队仓库建立联系

```shell script
git remote add upstream https://github.com/dgut-group-ten/nginx-proxy.git
```

工作前后要运行这几条命令,和团队项目保持同步

```shell script
git fetch upstream
git merge upstream/master
```

## 参考资料

- [Convert a float64 to an int in Go](https://stackoverflow.com/questions/8022389/convert-a-float64-to-an-int-in-go/8022789)
- [Go实战--golang中使用redis(redigo和go-redis/redis)](https://blog.csdn.net/wangshubo1989/article/details/75050024)
- [Parsing date/time strings which are not 'standard' formats](https://stackoverflow.com/questions/14106541/parsing-date-time-strings-which-are-not-standard-formats/14106561)
- [for example in the case of AUTH](https://gowalker.org/github.com/mediocregopher/radix.v2/pool)
- [删除redis所有KEY](https://blog.csdn.net/iloveyin/article/details/46813427)