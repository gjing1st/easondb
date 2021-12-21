# easondb

#### 介绍
1. go语言开发内存数据库(类似redis)
2. 参考数据库：nedb  
   github: https://github.com/louischatriot/nedb  
   简介：https://www.w3cschool.cn/nedbintro/nedbintro-eqsm27mb.html
3. 参考论文：https://riak.com/assets/bitcask-intro.pdf

#### 软件架构
1.  持久化时以字节形式存储。每个实体前保存该实体的头信息(key,value大小),读取时将根据头信息进行重新实例化实体
2.  字符串类型使用跳表做索引
3.  暂时使用默认配置，后续可根据配置文件进行配置。数据库存放路径，可通过config.DbPath进行修改，其他配置请查看config.go
4.  less is more 



#### 安装教程
```shell 
go get -u -v github.com/gjing1st/easondb
```
推荐使用 `go.mod`:
```
require github.com/gjing1st/easondb latest
```
#### 使用说明
参考根目录下string_test.go文件进行数据库初始化以及读写操作
入口可以从string Set查看


