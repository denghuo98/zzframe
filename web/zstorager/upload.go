package zstorager

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zcache"
	"github.com/denghuo98/zzframe/web/zcontext"
	"github.com/denghuo98/zzframe/web/zutils"
	"github.com/denghuo98/zzframe/web/zvalidate"
	"github.com/denghuo98/zzframe/zconsts"
)

// UploadDrive 存储驱动
type UploadDrive interface {
	// Upload 上传
	Upload(ctx context.Context, file *ghttp.UploadFile) (fullPath string, err error)
	// CreateMultipart 创建分片事件
	CreateMultipart(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error)
	// UploadPart 上传分片
	UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error)
}

// New 初始化存储驱动
func New(name ...string) UploadDrive {
	var (
		driveType = zconsts.UploadDriveLocal
		drive     UploadDrive
	)

	if len(name) > 0 && name[0] != "" {
		driveType = name[0]
	}

	switch driveType {
	case zconsts.UploadDriveLocal:
		drive = &LocalDrive{}
	// case consts.UploadDriveUCloud:
	// 	drive = &UCloudDrive{}
	// case consts.UploadDriveCos:
	// 	drive = &CosDrive{}
	// case consts.UploadDriveOss:
	// 	drive = &OssDrive{}
	// case consts.UploadDriveQiNiu:
	// 	drive = &QiNiuDrive{}
	// case consts.UploadDriveMinio:
	// 	drive = &MinioDrive{}
	default:
		panic(fmt.Sprintf("暂不支持的存储驱动:%v", driveType))
	}
	return drive
}

// DoUpload 上传入口
func DoUpload(ctx context.Context, typ string, file *ghttp.UploadFile) (result *entity.SysAttachment, err error) {
	if file == nil {
		err = gerror.New("文件必须!")
		return
	}

	meta, err := GetFileMeta(file)
	if err != nil {
		return
	}

	if err = ValidateFileMeta(typ, meta); err != nil {
		return
	}

	result, err = HasFile(ctx, meta.Md5)
	if err != nil {
		return
	}

	// 相同存储相同身份才复用
	if result != nil && result.Drive == config.Drive && result.MemberId == zcontext.GetUserId(ctx) {
		return
	}

	// 上传到驱动
	fullPath, err := New(config.Drive).Upload(ctx, file)
	if err != nil {
		return
	}
	// 写入附件记录
	return write(ctx, meta, fullPath)
}

// ValidateFileMeta 验证文件元数据
func ValidateFileMeta(typ string, meta *FileMeta) (err error) {
	switch typ {
	case KindImg:
		if !IsImgType(meta.Ext) {
			err = gerror.New("上传的文件不是图片")
			return
		}
		if config.ImageSize > 0 && meta.Size > config.ImageSize*1024*1024 {
			err = gerror.Newf("图片大小不能超过%vMB", config.ImageSize)
			return
		}

		if len(config.ImageType) > 0 && !lo.Contains(strings.Split(config.ImageType, `,`), meta.Ext) {
			err = gerror.New("上传图片类型未经允许")
			return
		}
	case KindDoc:
		if !IsDocType(meta.Ext) {
			err = gerror.New("上传的文件不是文档")
			return
		}
	case KindAudio:
		if !IsAudioType(meta.Ext) {
			err = gerror.New("上传的文件不是音频")
			return
		}
	case KindVideo:
		if !IsVideoType(meta.Ext) {
			err = gerror.New("上传的文件不是视频")
			return
		}
	case KindZip:
		if !IsZipType(meta.Ext) {
			err = gerror.New("上传的文件不是压缩文件")
			return
		}
	case KindOther:
		fallthrough
	default:
		// 默认为通用的文件上传
		if config.FileSize > 0 && meta.Size > config.FileSize*1024*1024 {
			err = gerror.Newf("文件大小不能超过%vMB", config.FileSize)
			return
		}

		if len(config.FileType) > 0 && !lo.Contains(strings.Split(config.FileType, `,`), meta.Ext) {
			err = gerror.Newf("上传文件类型未经允许:%v", meta.Ext)
			return
		}
	}
	return
}

// LastUrl 根据驱动获取最终文件访问地址
func LastUrl(ctx context.Context, fullPath, drive string) string {
	if zvalidate.IsURL(fullPath) {
		return fullPath
	}

	switch drive {
	case zconsts.UploadDriveLocal:
		return zutils.GetAddr(ctx) + "/" + fullPath
	case zconsts.UploadDriveUCloud:
		return config.UCloudEndpoint + "/" + fullPath
	case zconsts.UploadDriveCos:
		return config.CosBucketURL + "/" + fullPath
	case zconsts.UploadDriveOss:
		return config.OssBucketURL + "/" + fullPath
	case zconsts.UploadDriveQiNiu:
		return config.QiNiuDomain + "/" + fullPath
	case zconsts.UploadDriveMinio:
		return fmt.Sprintf("%s/%s/%s", config.MinioDomain, config.MinioBucket, fullPath)
	default:
		return fullPath
	}
}

// GetFileMeta 获取上传文件元数据
func GetFileMeta(file *ghttp.UploadFile) (meta *FileMeta, err error) {
	meta = new(FileMeta)
	meta.Filename = file.Filename
	meta.Size = file.Size
	meta.Ext = Ext(file.Filename)
	meta.Kind = GetFileKind(meta.Ext)
	meta.MimeType = GetFileMimeType(meta.Ext)

	// 兼容naiveUI
	naiveType := meta.MimeType
	if len(naiveType) == 0 {
		naiveType = "text/plain"
	}
	meta.NaiveType = naiveType

	// 计算md5值
	meta.Md5, err = CalcFileMd5(file)
	return
}

// GenFullPath 根据目录和文件类型生成一个绝对地址
func GenFullPath(basePath, ext string) string {
	fileName := strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6)
	fileName = fileName + ext
	return basePath + gtime.Date() + "/" + strings.ToLower(fileName)
}

// write 写入附件记录
func write(ctx context.Context, meta *FileMeta, fullPath string) (models *entity.SysAttachment, err error) {
	models = &entity.SysAttachment{
		AppId:     "admin",
		MemberId:  zcontext.GetUserId(ctx),
		Drive:     config.Drive,
		Size:      meta.Size,
		Path:      fullPath,
		FileUrl:   fullPath,
		Name:      meta.Filename,
		Kind:      meta.Kind,
		MimeType:  meta.MimeType,
		NaiveType: meta.NaiveType,
		Ext:       meta.Ext,
		Md5:       meta.Md5,
		Status:    zconsts.StatusEnabled,
	}

	id, err := GetModel(ctx).Data(models).OmitEmptyData().InsertAndGetId()
	if err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}
	models.Id = id
	return
}

// HasFile 检查附件是否存在
func HasFile(ctx context.Context, md5 string) (res *entity.SysAttachment, err error) {
	if err = GetModel(ctx).Where("md5", md5).Scan(&res); err != nil {
		err = gerror.Wrap(err, "检查文件hash时出现错误")
		return
	}

	if res == nil {
		return
	}

	// 只有在上传时才会检查md5值，如果附件存在则更新最后上传时间，保证上传列表更新显示在最前面
	if res.Id > 0 {
		update := g.Map{
			"status":     zconsts.StatusEnabled,
			"updated_at": gtime.Now(),
		}
		_, _ = GetModel(ctx).WherePri(res.Id).Data(update).Update()
	}
	return
}

// CheckMultipart 检查文件分片
func CheckMultipart(ctx context.Context, in *CheckMultipartParams) (res *CheckMultipartModel, err error) {
	res = new(CheckMultipartModel)

	meta := new(FileMeta)
	meta.Filename = in.FileName
	meta.Size = in.Size
	meta.Ext = Ext(in.FileName)
	meta.Kind = GetFileKind(meta.Ext)
	meta.MimeType = GetFileMimeType(meta.Ext)

	// 兼容naiveUI
	naiveType := "text/plain"
	if IsImgType(Ext(in.FileName)) {
		naiveType = ""
	}
	meta.NaiveType = naiveType
	meta.Md5 = in.Md5

	if err = ValidateFileMeta(in.UploadType, meta); err != nil {
		return
	}

	result, err := HasFile(ctx, in.Md5)
	if err != nil {
		return nil, err
	}

	// 文件已存在，直接返回。相同存储相同身份才复用
	if result != nil && result.Drive == config.Drive && result.MemberId == zcontext.GetUserId(ctx) {
		res.Attachment = result
		return
	}

	for i := 0; i < in.ShardCount; i++ {
		res.WaitUploadIndex = append(res.WaitUploadIndex, i+1)
	}

	in.meta = meta
	progress, err := GetOrCreateMultipartProgress(ctx, in)
	if err != nil {
		return nil, err
	}

	if len(progress.UploadedIndex) > 0 {
		res.WaitUploadIndex, _ = lo.Difference(progress.UploadedIndex, res.WaitUploadIndex)
	}

	if len(res.WaitUploadIndex) == 0 {
		res.WaitUploadIndex = make([]int, 0)
	}
	res.UploadId = progress.UploadId
	res.Progress = CalcUploadProgress(progress.UploadedIndex, progress.ShardCount)
	res.SizeFormat = humanize.IBytes(gconv.Uint64(progress.Meta.Size))
	return
}

// CalcUploadProgress 计算上传进度
func CalcUploadProgress(uploadedIndex []int, shardCount int) float64 {
	return gconv.Float64(len(uploadedIndex)) / gconv.Float64(shardCount) * 100
}

// GenUploadId 生成上传ID
func GenUploadId(ctx context.Context, md5 string) string {
	return fmt.Sprintf("%v:%v:%v@%v", md5, zcontext.GetUserId(ctx), "admin", config.Drive)
}

// GetOrCreateMultipartProgress 获取或创建分片上传事件进度
func GetOrCreateMultipartProgress(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error) {
	uploadId := GenUploadId(ctx, in.Md5)
	res, err = GetMultipartProgress(ctx, uploadId)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	return New(config.Drive).CreateMultipart(ctx, in)
}

// GetMultipartProgress 获取分片上传事件进度
func GetMultipartProgress(ctx context.Context, uploadId string) (res *MultipartProgress, err error) {
	key := fmt.Sprintf("%v:%v", zconsts.CacheMultipartUpload, uploadId)
	get, err := zcache.Instance().Get(ctx, key)
	if err != nil {
		return nil, err
	}
	err = get.Scan(&res)
	return
}

// CreateMultipartProgress 创建分片上传事件进度
func CreateMultipartProgress(ctx context.Context, in *MultipartProgress) (err error) {
	key := fmt.Sprintf("%v:%v", zconsts.CacheMultipartUpload, in.UploadId)
	return zcache.Instance().Set(ctx, key, in, time.Hour*24*7)
}

// UpdateMultipartProgress 更新分片上传事件进度
func UpdateMultipartProgress(ctx context.Context, in *MultipartProgress) (err error) {
	key := fmt.Sprintf("%v:%v", zconsts.CacheMultipartUpload, in.UploadId)
	return zcache.Instance().Set(ctx, key, in, time.Hour*24*7)
}

// DelMultipartProgress 删除分片上传事件进度
func DelMultipartProgress(ctx context.Context, in *MultipartProgress) (err error) {
	key := fmt.Sprintf("%v:%v", zconsts.CacheMultipartUpload, in.UploadId)
	_, err = zcache.Instance().Remove(ctx, key)
	return
}

// UploadPart 上传分片
func UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error) {
	in.mp, err = GetMultipartProgress(ctx, in.UploadId)
	if err != nil {
		return nil, err
	}
	if in.mp == nil {
		err = gerror.New("分片事件不存在，请重新上传！")
		return
	}

	if lo.Contains(in.mp.UploadedIndex, in.Index) {
		err = gerror.New("该分片已上传过了")
		return
	}

	res, err = New(config.Drive).UploadPart(ctx, in)
	if err != nil {
		return nil, err
	}
	return
}
