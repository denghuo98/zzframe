// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRole is the golang structure for table admin_role.
type AdminRole struct {
	Id        int64       `json:"id"        orm:"id"         description:"角色ID"`
	Name      string      `json:"name"      orm:"name"       description:"角色名称"`
	Key       string      `json:"key"       orm:"key"        description:"角色权限字符串"`
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`
	Sort      int         `json:"sort"      orm:"sort"       description:"排序"`
	Status    int         `json:"status"    orm:"status"     description:"角色状态"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`
}
