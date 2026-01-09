package zqueue

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/denghuo98/zzframe/web/zqueue/disk"
)

// MqProducer 消息生产者
type MqProducer interface {
	SendMsg(topic string, body string) (mqMsg MqMsg, err error)
	SendByteMsg(topic string, body []byte) (mqMsg MqMsg, err error)
	SendDelayMsg(topic string, body string, delay int64) (mqMsg MqMsg, err error)
}

// MqConsumer 消息消费者
type MqConsumer interface {
	ListenReceiveMsgDo(topic string, receiveDo func(mqMsg MqMsg)) (err error)
}

// MqMsg 消息体
type MqMsg struct {
	RunType   int       `json:"run_type"`
	Topic     string    `json:"topic"`
	MsgId     string    `json:"msg_id"`
	Offset    int64     `json:"offset"`
	Partition int32     `json:"partition"`
	Timestamp time.Time `json:"timestamp"`
	Body      []byte    `json:"body"`
}

const (
	_ = iota
	SendMsg
	ReceiveMsg
)

// Config 配置
type Config struct {
	Switch    bool   `json:"switch"`
	Driver    string `json:"driver"`
	GroupName string `json:"groupName"`
	Disk      *disk.Config
}

var (
	ctx                   = gctx.GetInitCtx()
	mqProducerInstanceMap map[string]MqProducer
	mqConsumerInstanceMap map[string]MqConsumer
	mutex                 sync.Mutex
	config                Config
)

func init() {
	mqProducerInstanceMap = make(map[string]MqProducer)
	mqConsumerInstanceMap = make(map[string]MqConsumer)
	configContent := g.Cfg().MustGet(ctx, "queue")
	if configContent != nil {
		if err := configContent.Struct(&config); err != nil {
			Logger().Warningf(ctx, "queue init err:%+v", err)
		}
	} else {
		g.Log().Warningf(ctx, "消息队列配置为空，使用默认配置")
		config = Config{
			Switch:    true,
			Driver:    "disk",
			GroupName: "default",
			Disk: &disk.Config{
				Path:         "./tmp/diskqueue",
				BatchSize:    100,
				BatchTime:    1,
				SegmentSize:  10485760,
				SegmentLimit: 3000,
			},
		}
	}
}

// InstanceConsumer 实例化消费者
func InstanceConsumer() (mqClient MqConsumer, err error) {
	return NewConsumer(config.GroupName)
}

// InstanceProducer 实例化生产者
func InstanceProducer() (mqClient MqProducer, err error) {
	return NewProducer(config.GroupName)
}

// NewProducer 初始化生产者实例
func NewProducer(groupName string) (mqClient MqProducer, err error) {
	if item, ok := mqProducerInstanceMap[groupName]; ok {
		return item, nil
	}

	if groupName == "" {
		err = gerror.New("mq groupName is empty.")
		return
	}

	switch config.Driver {

	case "disk":
		config.Disk.GroupName = groupName
		mqClient, err = RegisterDiskMqProducer(config.Disk)
	default:
		err = gerror.New("queue driver is not support")
	}

	if err != nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	mqProducerInstanceMap[groupName] = mqClient
	return
}

// NewConsumer 初始化消费者实例
func NewConsumer(groupName string) (mqClient MqConsumer, err error) {
	if groupName == "" {
		err = gerror.New("mq groupName is empty.")
		return
	}

	switch config.Driver {
	case "disk":
		config.Disk.GroupName = groupName
		mqClient, err = RegisterDiskMqConsumer(config.Disk)
	default:
		err = gerror.New("queue driver is not support")
	}

	if err != nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	mqConsumerInstanceMap[groupName] = mqClient
	return
}

// BodyString 返回消息体
func (m *MqMsg) BodyString() string {
	return string(m.Body)
}
