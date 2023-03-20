# Chat

#### 介绍
本项目参照微信功能，从0搭建一个web聊天系统。
#### 服务端实现以下功能：

- 一对一的文本消息传输，文件消息传输
- 一对多的文本消息传输，文件消息传输（群聊）
- 存储离线消息
  - 消息不丢失，可靠性，顺序性，每个消息有“已发送”/“已送达”/“已读”回执
- 支持用户登录，注册
- 支持好友关系
- 朋友圈功能，发表动态，评论功能
- 直播
- .......

#### 技术选型

开发语言Go

web框架：Gin

网络协议：http,websocket,webRTC

数据库：MySQL,Redis,MongoDB

消息队列：RabbitQM

分布式架构：GRPC,Protobuf

部署：Docker

#### 项目目录
│  .gitignore
│  docker-compose.yml
│  go.mod
│  go.sum
│  LICENSE
│  main.go
│      
├─api
│      Friends.go
│      Groups.go
│      Message.go
│      User.go
│      
├─config
│      config.go
│      config.ini
│      
├─grpc
│  │  rpcserver.go
│  │  
│  └─protofiles
├─Message
│      Message.go
│      Type.go
│      
├─middleware
│      JWT.go
│      
├─mongo
│      Document.go
│      mongo.go
│      
├─mysql
│      Friend.go
│      Group.go
│      Message.go
│      mysql.go
│      User.go
│      
├─rabbitmq
│      mq.go
│      
├─redis
│      redis.go
│      
├─router
│      Init.go
│      
├─service
│      FriendService.go
│      GroupServcie.go
│      MessageService.go
│      UserService.go
│      
├─test
│      main.go
│      
├─utils
│      errMsg.go
│      
└─webSocket
        UserConn.go
        webSocket.go

#### 注意
这里的mq消费我没有独立出去单独做一个程序，这里只是先进行测试，后续再做成一个独立的程序
   

#### 想法与实现
使用token进行用户鉴权

使用mongodb存储已发送但未读的消息

使用mysql存储已读消息、用户信息、群组信息、好友关系。

使用rabbitmq进行系统解耦，将消息持久化任务丢给mq

聊天消息功能使用websocket协议，方便服务器推送消息

其他用户请求使用http协议

文件传输功能我准备独立出来做成一个微服务。



#### 已经实现的功能
用户登录注册、用户鉴权

好友添加

群聊创建

webSocket连接、websocket消息传递

单聊功能、群聊功能、历史消息、消息持久化

消息基本有序，消息不丢失

#### 待实现功能
文件传输

论坛模块（类似pyq）

评论模块

直播模块


#### 前端部分受时间影响，尚未编写。
