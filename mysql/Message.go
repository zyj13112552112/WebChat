package mysql

import "chat/Message"

type MESSAGE Message.Message

func (msg *MESSAGE) INSERT() {
	DB.Model(&MESSAGE{}).Create(msg)
}
