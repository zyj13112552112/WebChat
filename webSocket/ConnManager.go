package webSocket

import (
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

var (
	ConnManager = make(map[string]*websocket.Conn)
	Register    = make(chan UserConn)
	UnRegister  = make(chan UserConn)
)

// Writer 写信息给客户端
func (uc *UserConn) Writer() {
	//todo 将离线消息按序发送，从mongo取出消息
	go func() {
		msgs := mongo.GetByToName(uc.Username)
		if len(msgs) == 0 {
			return
		}
		QuickSort(msgs, 0, len(msgs)-1)
		for _, v := range msgs {
			uc.WS.WriteJSON(v)
			bytes, _ := json.Marshal(ReadMessage(v))
			rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        bytes,
			})
			v.DELETE()
		}
	}()
	//发送好友请求
	go func() {
		requests := mongo.GetRequestByName(uc.Username)
		for _, v := range requests {
			uc.WS.WriteJSON(v)
			v.Delete()
		}
	}()
	//发送好友请求恢复
	go func() {
		responses := mongo.GetResponseByName(uc.Username)
		for _, v := range responses {
			uc.WS.WriteJSON(v)
			v.Delete()
		}
	}()
}

// QuickSort 消息按时间快速排序
func QuickSort(message []mongo.NoReadMessage, l, r int) {
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

// ReadMessage 未读消息转化为已读消息
func ReadMessage(noRead mongo.NoReadMessage) mysql.ReadMessage {
	return mysql.ReadMessage{
		From:    noRead.From,
		To:      noRead.To,
		Body:    noRead.Body,
		Time:    noRead.Time,
		Groupid: noRead.GroupID,
	}
}

// 读取客户端信息
func (uc *UserConn) Read() {
	for {
		var Msg mongo.NoReadMessage
		err := uc.WS.ReadJSON(&Msg)
		if err != nil {
			UnRegister <- *uc
			return
		}
		if Msg.GroupID == "" { //单聊
			//todo
			//对方用户在线，直接发送，需要保证数据有序性
			if ConnManager[Msg.To] != nil {
				ConnManager[Msg.To].WriteJSON(Msg)
				bytes, _ := json.Marshal(ReadMessage(Msg))
				rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        bytes,
				})
			} else { //不在线，丢进队列解耦，监听进行写数据库
				offLineMsg, _ := json.Marshal(Msg)
				//丢进队列进行写数据库持久化
				rabbitmq.CH.Publish("", rabbitmq.OFFLINE.Name, false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        offLineMsg,
				})
			}
		} else { //群聊
			//todo
			//遍历群聊用户，逐一发送
			usernames, _ := mysql.SelectByGroupName(Msg.GroupID)
			for _, v := range usernames {
				if v == Msg.From {
					continue
				}
				var offlineMsg mongo.NoReadMessage = mongo.NoReadMessage{
					GroupID: Msg.GroupID,
					From:    Msg.From,
					To:      v,
					Body:    Msg.Body,
					Time:    Msg.Time,
				}
				//对方用户在线，直接发送
				if ConnManager[v] != nil {
					ConnManager[v].WriteJSON(offlineMsg)
					bytes, _ := json.Marshal(ReadMessage(offlineMsg))
					rabbitmq.CH.Publish("", rabbitmq.READQUEUE.Name, false, false, amqp.Publishing{
						ContentType: "text/plain",
						Body:        bytes,
					})
				} else { //不在线，丢进队列解耦，进行写数据库
					offLineMsg, _ := json.Marshal(offlineMsg)
					rabbitmq.CH.Publish("", rabbitmq.OFFLINE.Name, false, false, amqp.Publishing{
						ContentType: "text/plain",
						Body:        offLineMsg,
					})
				}
			}
		}

	}

}
