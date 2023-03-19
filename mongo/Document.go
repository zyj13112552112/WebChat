package mongo

import (
	"chat/Message"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// AddFriendRequestDocument 好友请求消息
type AddFriendRequestDocument Message.Message

func (document *AddFriendRequestDocument) Insert() {
	CollectionAddRequest.InsertOne(context.TODO(), document)
}
func (document *AddFriendRequestDocument) Delete() {
	filter := bson.D{
		{"usersend", document.UserSend},
		{"useracc", document.UserAcc},
		{"groupid", document.GroupId},
		{"body", document.Body},
		{"time", document.Time},
	}
	CollectionAddRequest.DeleteOne(context.TODO(), filter)
}
func GetAddFriendRequestDocument(accUser string) (res []AddFriendRequestDocument) {
	filter := bson.D{
		{"useracc", accUser},
	}
	rows, _ := CollectionAddRequest.Find(context.TODO(), filter)
	for rows.Next(context.TODO()) {
		var msg AddFriendRequestDocument
		rows.Decode(&msg)
		res = append(res, msg)
	}
	rows.Close(context.TODO())
	return
}

// AddFriendResponseDocument 好友请求回复消息
type AddFriendResponseDocument Message.Message

func (document *AddFriendResponseDocument) Insert() {
	CollectionAddResponse.InsertOne(context.TODO(), document)
}
func (document *AddFriendResponseDocument) Delete() {
	filter := bson.D{
		{"usersend", document.UserSend},
		{"useracc", document.UserAcc},
		{"groupid", document.GroupId},
		{"body", document.Body},
		{"time", document.Time},
	}
	CollectionAddResponse.DeleteOne(context.TODO(), filter)
}
func GetAddFriendResponseDocument(accUser string) (res []AddFriendResponseDocument) {
	filter := bson.D{
		{"useracc", accUser},
	}
	rows, _ := CollectionAddResponse.Find(context.TODO(), filter)
	for rows.Next(context.TODO()) {
		var msg AddFriendResponseDocument
		rows.Decode(&msg)
		res = append(res, msg)
	}
	rows.Close(context.TODO())
	return
}

// ChatMessage 聊天消息
type ChatMessage Message.Message

func (document *ChatMessage) Insert() {
	CollectionMsg.InsertOne(context.TODO(), document)
}
func (document *ChatMessage) Delete() {
	filter := bson.D{
		{"usersend", document.UserSend},
		{"useracc", document.UserAcc},
		{"groupid", document.GroupId},
		{"body", document.Body},
		{"time", document.Time},
	}
	CollectionMsg.DeleteOne(context.TODO(), filter)
}
func GetChatMessage(accUser string) (res []ChatMessage) {
	filter := bson.D{
		{"useracc", accUser},
	}
	rows, _ := CollectionMsg.Find(context.TODO(), filter)
	for rows.Next(context.TODO()) {
		var msg ChatMessage
		rows.Decode(&msg)
		res = append(res, msg)
	}
	rows.Close(context.TODO())
	if len(res) != 0 {
		QuickSort(res, 0, len(res)-1)
	}
	return
}

// QuickSort 快速排序,保证消息有序
func QuickSort(message []ChatMessage, l, r int) {
	ll, rr := l, r
	var mid int64 = message[(l+r)/2].Time
	for {
		if ll >= rr {
			break
		}
		for {
			if message[ll].Time < mid {
				ll++
			} else {
				break
			}
		}
		for {
			if message[rr].Time > mid {
				rr--
			} else {
				break
			}
		}
		if ll <= rr {
			message[ll], message[rr] = message[rr], message[ll]
			ll++
			rr--
		}
	}
	if ll < r {
		QuickSort(message, ll, r)
	}
	if rr > l {
		QuickSort(message, l, rr)
	}
}
