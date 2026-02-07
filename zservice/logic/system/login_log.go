package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zlocation"
	"github.com/denghuo98/zzframe/web/zqueue"
	"github.com/denghuo98/zzframe/web/zutils"
	"github.com/denghuo98/zzframe/zconsts"
	commonSchema "github.com/denghuo98/zzframe/zschema/common"
	systemSchema "github.com/denghuo98/zzframe/zschema/system"
	"github.com/denghuo98/zzframe/zservice"
)

func init() {
	zservice.RegisterSysLoginLog(NewSysLoginLog())
}

type sSysLoginLog struct{}

func NewSysLoginLog() *sSysLoginLog {
	return &sSysLoginLog{}
}

// Push 推送登录日志（通常异步推送）
func (s *sSysLoginLog) Push(ctx g.Ctx, in *systemSchema.SysLoginLogPushInput) (err error) {
	if in.Response == nil {
		in.Response = new(commonSchema.SiteLoginOutput)
	}

	r := g.RequestFromCtx(ctx)
	if r == nil {
		g.Log().Warningf(ctx, "request from ctx is nil")
		return
	}

	clientIp := zlocation.GetClientIp(r)
	ipData, err := zlocation.GetLocation(ctx, clientIp)
	if err != nil {
		g.Log().Warningf(ctx, "get location fail, ip:%v, err:%v", clientIp, err)
		return
	}

	if ipData == nil {
		ipData = new(zlocation.IpLocationData)
	}

	var models entity.SysLoginLog
	models.ReqId = gctx.CtxId(ctx)
	models.MemberId = in.Response.Id
	models.Username = in.Response.Username
	models.LoginAt = gtime.Now()
	models.LoginIp = clientIp
	models.UserAgent = r.UserAgent()
	models.Province = ipData.Province
	models.City = ipData.City
	models.Status = zconsts.StatusEnabled

	if in.Error != nil {
		models.Status = zconsts.StatusDisable
		models.ErrMsg = in.Error.Error()
	}

	models.Response = gjson.New(zconsts.NilJsonToString)
	if in.Response != nil {
		models.Response = gjson.New(in.Response)
	}

	if err = zqueue.Push(zconsts.QueueLoginLogTopic, models); err != nil {
		g.Log().Warningf(ctx, "push login log to queue failed, err:%+v, data:%+v", err, models)
	}
	return
}

func (s *sSysLoginLog) RealWrite(ctx g.Ctx, data entity.SysLoginLog) (err error) {
	_, err = dao.SysLoginLog.Ctx(ctx).Data(data).OmitEmpty().Insert()
	return
}

func (s *sSysLoginLog) Delete(ctx g.Ctx, in *systemSchema.SysLoginLogDeleteInput) (err error) {
	if in.Id <= 0 {
		return gerror.New("ID不能为空")
	}
	_, err = dao.SysLoginLog.Ctx(ctx).WherePri(in.Id).Delete()
	return
}

func (s *sSysLoginLog) List(ctx g.Ctx, in *systemSchema.SysLoginLogListInput) (out *systemSchema.SysLoginLogListOutput, totalCount int, err error) {
	var (
		m    = dao.SysLoginLog.Ctx(ctx)
		cols = dao.SysLoginLog.Columns()
	)

	// 筛选状态
	if in.Status > 0 {
		m = m.Where(cols.Status, in.Status)
	}

	// 筛选登录时间
	if len(in.LoginAt) == 2 {
		m = m.WhereBetween(cols.LoginAt, in.LoginAt[0], in.LoginAt[1])
	}

	// 筛选登录IP
	if in.LoginIp != "" {
		m = m.Where(cols.LoginIp, in.LoginIp)
	}

	// 筛选账号
	if in.Username != "" {
		m = m.Where(cols.Username, in.Username)
	}

	totalCount, err = m.Count()
	if err != nil || totalCount == 0 {
		return
	}

	var list []*systemSchema.SysLoginLogListOuputItem
	if err = m.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, zconsts.ErrorORM)
	}

	for _, item := range list {
		item.Os = zutils.GetOs(item.UserAgent)
		item.Browser = zutils.GetBrowser(item.UserAgent)
	}
	out = new(systemSchema.SysLoginLogListOutput)
	out.List = list
	return
}

// Export 导出登录日志为 Excel
func (s *sSysLoginLog) Export(ctx g.Ctx, in *systemSchema.SysLoginLogExportInput) (filePath string, err error) {
	const maxExportCount = 10000 // 最大导出数量限制

	var (
		m    = dao.SysLoginLog.Ctx(ctx)
		cols = dao.SysLoginLog.Columns()
	)

	// 筛选状态
	if in.Status > 0 {
		m = m.Where(cols.Status, in.Status)
	}

	// 筛选登录时间
	if len(in.LoginAt) == 2 {
		m = m.WhereBetween(cols.LoginAt, in.LoginAt[0], in.LoginAt[1])
	}

	// 筛选登录IP
	if in.LoginIp != "" {
		m = m.Where(cols.LoginIp, in.LoginIp)
	}

	// 筛选账号
	if in.Username != "" {
		m = m.Where(cols.Username, in.Username)
	}

	// 检查数据量
	totalCount, err := m.Count()
	if err != nil {
		return "", gerror.Wrap(err, "查询数据量失败")
	}

	if totalCount > maxExportCount {
		return "", gerror.Newf("导出数据量超过限制（%d条），当前查询结果有%d条。请设置筛选条件（如时间范围、账号、IP等）来减少数据量", maxExportCount, totalCount)
	}

	if totalCount == 0 {
		return "", gerror.New("没有可导出的数据")
	}

	// 查询所有数据
	var list []*systemSchema.SysLoginLogListOuputItem
	if err = m.OrderDesc(cols.Id).Scan(&list); err != nil {
		return "", gerror.Wrap(err, zconsts.ErrorORM)
	}

	// 处理操作系统和浏览器信息
	for _, item := range list {
		item.Os = zutils.GetOs(item.UserAgent)
		item.Browser = zutils.GetBrowser(item.UserAgent)
	}

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			g.Log().Errorf(ctx, "关闭 Excel 文件失败: %v", err)
		}
	}()

	sheetName := "登录日志"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", gerror.Wrap(err, "创建工作表失败")
	}

	// 删除默认的 Sheet1
	if err = f.DeleteSheet("Sheet1"); err != nil {
		return "", gerror.Wrap(err, "删除默认工作表失败")
	}

	// 设置活动工作表
	f.SetActiveSheet(index)

	// 设置表头
	headers := []string{"ID", "请求ID", "用户ID", "用户名", "登录时间", "登录IP", "省份", "城市", "操作系统", "浏览器", "用户代理", "错误信息", "状态", "创建时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err = f.SetCellValue(sheetName, cell, header); err != nil {
			return "", gerror.Wrapf(err, "设置表头失败: %s", header)
		}
	}

	// 设置表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return "", gerror.Wrap(err, "创建表头样式失败")
	}

	// 应用表头样式
	for i := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err = f.SetCellStyle(sheetName, cell, cell, headerStyle); err != nil {
			return "", gerror.Wrapf(err, "应用表头样式失败: %s", cell)
		}
	}

	// 写入数据
	for rowIdx, item := range list {
		row := rowIdx + 2 // 从第2行开始（第1行是表头）

		// 状态文本
		statusText := "未知"
		if item.Status == zconsts.StatusEnabled {
			statusText = "成功"
		} else if item.Status == zconsts.StatusDisable {
			statusText = "失败"
		}

		// 登录时间
		loginAtStr := ""
		if item.LoginAt != nil {
			loginAtStr = item.LoginAt.Format("Y-m-d H:i:s")
		}

		// 创建时间
		createdAtStr := ""
		if item.CreatedAt != nil {
			createdAtStr = item.CreatedAt.Format("Y-m-d H:i:s")
		}

		// 写入数据行
		rowData := []interface{}{
			item.Id,
			item.ReqId,
			item.MemberId,
			item.Username,
			loginAtStr,
			item.LoginIp,
			item.Province,
			item.City,
			item.Os,
			item.Browser,
			item.UserAgent,
			item.ErrMsg,
			statusText,
			createdAtStr,
		}

		for colIdx, value := range rowData {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			if err = f.SetCellValue(sheetName, cell, value); err != nil {
				return "", gerror.Wrapf(err, "写入数据失败: 行%d, 列%d", row, colIdx+1)
			}
		}
	}

	// 自动调整列宽
	for i := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		if err = f.SetColWidth(sheetName, colName, colName, 15); err != nil {
			g.Log().Warningf(ctx, "设置列宽失败: %s, %v", colName, err)
		}
	}

	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "excel_export")
	if !gfile.Exists(tempDir) {
		if err = gfile.Mkdir(tempDir); err != nil {
			return "", gerror.Wrap(err, "创建临时目录失败")
		}
	}

	// 生成文件名
	fileName := fmt.Sprintf("login_log_%s.xlsx", gtime.Now().Format("YmdHis"))
	filePath = filepath.Join(tempDir, fileName)

	// 保存文件
	if err = f.SaveAs(filePath); err != nil {
		return "", gerror.Wrap(err, "保存 Excel 文件失败")
	}

	return filePath, nil
}
