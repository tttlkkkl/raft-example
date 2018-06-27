# raft-example
raft 学习测试代码
- 启动
```golang
go run main.go
```
- 简单的 web 接口
```shell
#获取一个key的值
curl 127.0.0.1:6061?get=k
#设置一个key的值
curl 127.0.0.1:6061?set=k:v
#删除一个key
curl 127.0.0.1:6061?delete=k
```
