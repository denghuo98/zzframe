// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminMenu is the golang structure for table admin_menu.
type AdminMenu struct {
	Id          int64       `json:"id"          orm:"id"          description:"菜单ID"`
	Pid         int64       `json:"pid"         orm:"pid"         description:"父菜单ID"`
	Title       string      `json:"title"       orm:"title"       description:"菜单名称"`
	Name        string      `json:"name"        orm:"name"        description:"名称编码"`
	Path        string      `json:"path"        orm:"path"        description:"路由地址"`
	Icon        string      `json:"icon"        orm:"icon"        description:"菜单图标"`
	Type        int         `json:"type"        orm:"type"        description:"菜单类型（1目录 2菜单 3按钮）"`
	Redirect    string      `json:"redirect"    orm:"redirect"    description:"重定向地址"`
	Component   string      `json:"component"   orm:"component"   description:"组件路径"`
	IsFrame     int         `json:"isFrame"     orm:"is_frame"    description:"是否内嵌"`
	FrameSrc    string      `json:"frameSrc"    orm:"frame_src"   description:"内联外部地址"`
	Hidden      int         `json:"hidden"      orm:"hidden"      description:"是否隐藏"`
	Sort        int         `json:"sort"        orm:"sort"        description:"排序"`
	Remark      string      `json:"remark"      orm:"remark"      description:"备注"`
	Status      int         `json:"status"      orm:"status"      description:"菜单状态"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:"更新时间"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:"创建时间"`
	Permissions string      `json:"permissions" orm:"permissions" description:"菜单权限"`
}
