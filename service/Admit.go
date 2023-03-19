package service

import (
	"chat/mongo"
	"chat/mysql"
	"chat/rabbitmq"
	"encoding/json"
	"github.com/streadway/amqp"
)

// Admit 好友请求答复
type Admit struct {
	UserA  string `json:"usera" form:"usera" binding:"required"`
	UserB  string `json:"userb" form:"userb" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
}

func (admin *Admit) ADD() (code int, err error) {
	if admin.Status == "true" {
		//同意
		var f mysql.Friend = mysql.Friend{
			UsernameAB: admin.UserA,
			UsernameBA: admin.UserB,
		}
		mysql.AddFriend(f)
	}
	//发送结果给请求添加好友方
	var M mongo.ADDresponse = mongo.ADDresponse{
		UserA:  admin.UserA,
		UserB:  admin.UserB,
		Status: admin.Status,
	}
	marshal, _ := json.Marshal(M)
	rabbitmq.CH.Publish("", rabbitmq.ADDRESPONSE.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        marshal,
	})
	return
}
