// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package dto

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Reply is the golang structure of table gf_reply for DAO operations like Where/Data.
type Reply struct {
	g.Meta     `orm:"table:gf_reply, dto:true"`
	Id         interface{} // 回复ID
	ParentId   interface{} // 回复对应的上一级回复ID(没有的话默认为0)
	Title      interface{} // 回复标题
	Content    interface{} // 回复内容
	TargetType interface{} // 评论类型: content, reply
	TargetId   interface{} // 对应内容ID，可能回复的是另一个回复，所以这里没有使用content_id
	UserId     interface{} // 网站用户ID
	ZanCount   interface{} // 赞
	CaiCount   interface{} // 踩
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time //
}
