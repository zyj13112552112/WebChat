package service

import (
	"chat/Message"
	"chat/mysql"
	"chat/rabbitmq"
	"chat/utils"
	"chat/webSocket"
	"encoding/json"
	"github.com/streadway/amqp"
	"time"
)

type FriendService Message.Message

func (fs *FriendService) AddRequest() (code int) {
	fs.Type = Message.AddFriendRequest
	fs.Time = time.Now().Unix()
	//转发好友请求
	accUser := fs.UserAcc
	//对方在线立即转发，不在线丢进消息队列
	if webSocket.ConnManager[accUser] != nil {
		webSocket.ConnManager[accUser].WriteJSON(fs)
		var msg mysql.MESSAGE = mysql.MESSAGE(*fs)
		msg.INSERT() //持久化到数据库
	} else {
		body, _ := json.Marshal(fs)
		rabbitmq.CH.Publish("", rabbitmq.ADDREQUEST.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	}
	return utils.SUCCESS
}

func (fs *FriendService) AddResponse() (code int) {
	fs.Type = Message.AddFriendResponse
	accUser := fs.UserAcc
	//同意，发送消息给请求添加好友的用户,同时进行mysql写入
	if fs.Agree == 0 {
		//写入数据库
		go func(fs *FriendService) {
			var friend mysql.Friend = mysql.Friend{
				UsernameAB: fs.UserSend,
				UsernameBA: fs.UserAcc,
			}
			friend.Insert()
		}(fs)
		//发送消息
		if webSocket.ConnManager[accUser] != nil {
			webSocket.ConnManager[accUser].WriteJSON(fs)
			var msg mysql.MESSAGE = mysql.MESSAGE(*fs)
			msg.INSERT() //持久化到数据库
		} else {
			body, _ := json.Marshal(fs)
			rabbitmq.CH.Publish("", rabbitmq.ADDRESPONSE.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})
		}
	}
	//不同意，发送消息给请求添加好友的用户
	if fs.Agree == 1 {
		if webSocket.ConnManager[accUser] != nil {
			webSocket.ConnManager[accUser].WriteJSON(fs)
			var msg mysql.MESSAGE = mysql.MESSAGE(*fs)
			msg.INSERT() //持久化到数据库
		} else {
			body, _ := json.Marshal(fs)
			rabbitmq.CH.Publish("", rabbitmq.ADDRESPONSE.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})
		}
	}
	return utils.SUCCESS
}
