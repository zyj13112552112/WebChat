package api

import (
	"chat/service"
	"chat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// AddFriendRequest 添加好友请求
// form表单提交
// userA , userB
func AddFriendRequest(c *gin.Context) {
	var fs service.FriendService
	err := c.ShouldBind(&fs)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(200, gin.H{
			"msg": err.Error(),
		})
		return
	}
	code := fs.AddRequest()
	c.JSON(200, gin.H{
		"msg": utils.GetErrMsg(code),
	})
}

// AddFriendResponse 是否同意好友请求
func AddFriendResponse(c *gin.Context) {
	var fs service.FriendService
	err := c.ShouldBind(&fs)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": err.Error(),
		})
		return
	}
	code := fs.AddResponse()
	c.JSON(200, gin.H{
		"msg": utils.GetErrMsg(code),
	})
}

func DeleteFriend(c *gin.Context) {

}
