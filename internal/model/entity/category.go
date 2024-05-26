// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Category is the golang structure for table category.
type Category struct {
	Id          uint        `json:"id"          description:"分类ID，自增主键"`
	ContentType string      `json:"contentType" description:"内容类型：topic, ask, article, reply"`
	Key         string      `json:"key"         description:"栏目唯一键名，用于程序部分场景硬编码，一般不会用得到"`
	ParentId    uint        `json:"parentId"    description:"父级分类ID，用于层级管理"`
	UserId      uint        `json:"userId"      description:"创建的用户ID"`
	Name        string      `json:"name"        description:"分类名称"`
	Sort        uint        `json:"sort"        description:"排序，数值越低越靠前，默认为添加时的时间戳，可用于置顶"`
	Thumb       string      `json:"thumb,omitempty"       description:"封面图"`
	Brief       string      `json:"brief,omitempty"       description:"简述"`
	Content     string      `json:"content,omitempty"     description:"详细介绍"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"修改时间"`
}
