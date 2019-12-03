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