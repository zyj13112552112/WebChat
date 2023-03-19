package webSocket

import (
	"chat/Message"
	"chat/mongo"
	"chat/mysql"
	"chat/rabbitmq"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type UserConn struct {
	Username string
	WS       *websocket.Conn
}

// Writer 返回信息给客户端
func (uc *UserConn) Writer() {
	//todo 将离线消息按序发送，从mongo取出消息
	go func() {
		msgs := mongo.GetChatMessage(uc.Username)
		for _, msg := range msgs {
			go func(msg1 mongo.ChatMessage) {
				uc.WS.WriteJSON(msg1) //发送给客户端
				msg1.Delete()         //从mongo中删除
				bytes, _ := json.Marshal(msg1)
				//已读丢给消息队列进行持久化
				rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        bytes,
				})
			}(msg)
		}
	}()
	//todo 发送好友请求
	go func() {
		msgs := mongo.GetAddFriendRequestDocument(uc.Username)
		for _, msg := range msgs {
			go func(msg1 mongo.AddFriendRequestDocument) {
				uc.WS.WriteJSON(msg1) //发送给客户端
				msg1.Delete()         //从mongo中删除
				bytes, _ := json.Marshal(msg1)
				//已读丢给消息队列进行持久化
				rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        bytes,
				})
			}(msg)
		}
	}()
	//todo 发送好友请求回复
	go func() {
		msgs := mongo.GetAddFriendResponseDocument(uc.Username)
		for _, msg := range msgs {
			go func(msg1 mongo.AddFriendResponseDocument) {
				uc.WS.WriteJSON(msg1) //发送给客户端
				msg1.Delete()         //从mongo中删除
				bytes, _ := json.Marshal(msg1)
				//已读丢给消息队列进行持久化
				rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        bytes,
				})
			}(msg)
		}
	}()
}

// Read 读取客户端信息
func (uc *UserConn) Read() {
	for {
		var msg mongo.ChatMessage
		err := uc.WS.ReadJSON(&msg)
		if err != nil {
			UnRegister <- *uc
			break
		}
		if msg.Type == Message.ChatMessage {
			bytes, _ := json.Marshal(msg)
			if ConnManager[msg.UserAcc] != nil {
				ConnManager[msg.UserAcc].WriteJSON(msg)
				//已读丢给消息队列进行持久化
				rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        bytes,
				})
			} else {
				//对方不在线，持久化到mongo
				chatMsg := mongo.ChatMessage(msg)
				chatMsg.Insert()
			}
		} else if msg.Type == Message.GroupChatMessage {
			//查找群聊所有成员
			users, _ := mysql.SelectByGroupName(msg.GroupId)
			for _, user := range users {
				if user == msg.UserSend {
					continue
				}
				msg.UserAcc = user
				bytes, _ := json.Marshal(msg)
				if ConnManager[user] != nil {
					ConnManager[user].WriteJSON(msg)
					//已读丢给消息队列进行持久化
					rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
						ContentType: "text/plain",
						Body:        bytes,
					})
				} else {
					//对方不在线，丢给rabbitmq持久化到mongo
					rabbitmq.CH.Publish("", rabbitmq.OFFLINE.Name, false, false, amqp.Publishing{
						ContentType: "text/plain",
						Body:        bytes,
					})
				}
			}
		}
	}
}
