##### gonet go封装的tcp/websocket网络库

###### `go mod download github.com/cwloo/gonet@latest`

###### 1.`defer性能差禁用`
###### 2.`字符串连接用strings.joins`
###### 3.`用sync.Pool 避免频繁GC卡顿`
###### 4.`定期手动GC，避免内存暴涨`
###### 5.`slice内存泄漏避坑，goroutine泄漏避坑`
###### 6.`减少锁调用，能避免则避免或用原子锁`
###### 7.`使用goroutine池，频繁启动耗费cpu`
###### 8.`尽量避免多个goroutine竞争访问共享资源`



![image](https://github.com/cwloo/gonet/blob/main/res/log.png)


![image](https://github.com/cwloo/gonet/blob/main/res/server.png)