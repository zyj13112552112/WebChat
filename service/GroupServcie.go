package service

import (
	"chat/Model"
	"chat/utils"
)

type GroupService struct {
	GroupName      string   `json:"groupname" binding:"required,min=1"` //群名称
	CreateUserName string   `json:"createusername"`                     //创建者
	Members        []string `json:"members"`                            //成员
}

func (gs *GroupService) CreateGroup() (code int, groupid string) {
	groupid = Model.CreateGroup(gs.GroupName, gs.CreateUserName, gs.Members)
	if groupid == "" {
		return utils.GROUP_CREATE_ERR, groupid
	}
	return utils.SUCCESS, groupid
}
