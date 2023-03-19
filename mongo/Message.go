package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// NoReadMessage 为读数据持久化
// 前后端消息交互
type NoReadMessage struct {
	From    string //消息谁发送的
	To      string //消息发给谁
	GroupID string //指定群聊
	Body    []byte //消息内容
	Time    int64  //保证消息有序
}

func Insert(message NoReadMessage) {
	CollectionMsg.InsertOne(context.TODO(), message)
}

func GetByToName(To string) (res []NoReadMessage) {
	filter := bson.D{
		{"to", To},
	}
	rows, _ := CollectionMsg.Find(context.TODO(), filter)
	for rows.Next(context.TODO()) {
		var msgs NoReadMessage
		rows.Decode(&msgs)
		res = append(res, msgs)
	}
	rows.Close(context.TODO())
	return
}

func (msg *NoReadMessage) DELETE() {
	Filter := bson.D{
		{"from", msg.From},
		{"to", msg.To},
		{"groupid", msg.GroupID},
		{"body", msg.Body},
		{"time", msg.Time},
	}
	CollectionMsg.DeleteOne(context.TODO(), Filter)
}
