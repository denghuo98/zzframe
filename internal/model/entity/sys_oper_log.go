// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysOperLog is the golang structure for table sys_oper_log.
type SysOperLog struct {
	Id         int64       `json:"id"         orm:"id"          description:"日志ID"`
	Title      string      `json:"title"      orm:"title"       description:"操作标题"`
	Method     string      `json:"method"     orm:"method"      description:"请求方法"`
	RequestUri string      `json:"requestUri" orm:"request_uri" description:"请求URI"`
	OperParam  *gjson.Json `json:"operParam"  orm:"oper_param"  description:"请求参数"`
	Result     *gjson.Json `json:"result"     orm:"result"      description:"响应结果"`
	OperId     int64       `json:"operId"     orm:"oper_id"     description:"操作人ID"`
	OperName   string      `json:"operName"   orm:"oper_name"   description:"操作人名称"`
	OperIp     string      `json:"operIp"     orm:"oper_ip"     description:"操作IP"`
	BizId      string      `json:"bizId"      orm:"biz_id"      description:"业务ID"`
	BizType    string      `json:"bizType"    orm:"biz_type"    description:"业务类型"`
	Diff       *gjson.Json `json:"diff"       orm:"diff"        description:"数据变更差异"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"修改时间"`
}
