package mysql

import "github.com/jinzhu/gorm"

//历史消息

// ReadMessage 已读数据持久化
type ReadMessage struct {
	gorm.Model
	From    string `gorm:"from"`
	To      string `gorm:"to"`
	Body    []byte `gorm:"body"`
	Time    int64  `gorm:"time"`
	Groupid string `gorm:"groupid"`
}

func (readMsg ReadMessage) Insert() {
	DB.Model(&ReadMessage{}).Create(&readMsg)
}
