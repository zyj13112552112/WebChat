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

C:.
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
        


