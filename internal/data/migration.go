package data

import (
	"context"
	"log/slog"

	"github.com/gowvp/owl/internal/core/bz"
	"github.com/gowvp/owl/internal/core/ipc"
	"github.com/ixugo/goddd/domain/uniqueid"
	"github.com/ixugo/goddd/pkg/orm"
	"gorm.io/gorm"
)

// StreamPush 旧的推流模型（用于迁移）
type StreamPush struct {
	ID             string    `gorm:"primaryKey"`
	App            string    `gorm:"column:app"`
	Stream         string    `gorm:"column:stream"`
	IsAuthDisabled bool      `gorm:"column:is_auth_disabled"`
	Session        string    `gorm:"column:session"`
	Status         string    `gorm:"column:status"`
	PushedAt       *orm.Time `gorm:"column:pushed_at"`
	StoppedAt      *orm.Time `gorm:"column:stopped_at"`
	MediaServerID  string    `gorm:"column:media_server_id"`
	Name           string    `gorm:"column:name"`
	CreatedAt      orm.Time  `gorm:"column:created_at"`
	UpdatedAt      orm.Time  `gorm:"column:updated_at"`
}

func (*StreamPush) TableName() string {
	return "stream_pushs"
}

// StreamProxy 旧的代理模型（用于迁移）
type StreamProxy struct {
	ID                        string   `gorm:"primaryKey"`
	App                       string   `gorm:"column:app"`
	Stream                    string   `gorm:"column:stream"`
	Name                      string   `gorm:"column:name"`
	SourceURL                 string   `gorm:"column:source_url"`
	Transport                 int      `gorm:"column:transport"`
	TimeoutS                  int      `gorm:"column:timeout_s"`
	EnabledAudio              bool     `gorm:"column:enabled_audio"`
	EnabledRemoveNoneReader   bool     `gorm:"column:enabled_remove_none_reader"`
	EnabledDisabledNoneReader bool     `gorm:"column:enabled_disabled_none_reader"`
	StreamKey                 string   `gorm:"column:stream_key"`
	Pulling                   bool     `gorm:"column:pulling"`
	Enabled                   bool     `gorm:"column:enabled"`
	CreatedAt                 orm.Time `gorm:"column:created_at"`
	UpdatedAt                 orm.Time `gorm:"column:updated_at"`
}

func (*StreamProxy) TableName() string {
	return "stream_proxys"
}

// MigrateStreamData 迁移 stream_pushs 和 stream_proxys 数据到 channels 表
// 迁移完成后，旧表数据保留，建议手动确认后删除
func MigrateStreamData(db *gorm.DB, uni uniqueid.Core) error {
	ctx := context.Background()

	// 检查是否存在旧表
	hasStreamPushs := db.Migrator().HasTable("stream_pushs")
	hasStreamProxys := db.Migrator().HasTable("stream_proxys")

	if !hasStreamPushs && !hasStreamProxys {
		slog.Info("没有需要迁移的旧表数据")
		return nil
	}

	// 创建虚拟设备用于存放迁移的通道
	var rtmpDevice, rtspDevice *ipc.Device

	// RTMP 虚拟设备
	if hasStreamPushs {
		id := uni.UniqueID(bz.IDPrefixRTMP)
		rtmpDevice = &ipc.Device{
			ID:       id,
			DeviceID: id,
			Name:     "RTMP 迁移设备",
			Type:     ipc.TypeRTMP,
			IsOnline: true,
		}
		if err := db.WithContext(ctx).FirstOrCreate(rtmpDevice, "device_id = ?", rtmpDevice.DeviceID).Error; err != nil {
			slog.Error("创建 RTMP 虚拟设备失败", "err", err)
			return err
		}
		slog.Info("RTMP 虚拟设备已创建/存在", "id", rtmpDevice.ID)
	}

	// RTSP 虚拟设备
	if hasStreamProxys {
		id := uni.UniqueID(bz.IDPrefixRTSP)
		rtspDevice = &ipc.Device{
			ID:       id,
			DeviceID: id,
			Name:     "RTSP 迁移设备",
			Type:     ipc.TypeRTSP,
			IsOnline: true,
		}
		if err := db.WithContext(ctx).FirstOrCreate(rtspDevice, "device_id = ?", rtspDevice.DeviceID).Error; err != nil {
			slog.Error("创建 RTSP 虚拟设备失败", "err", err)
			return err
		}
		slog.Info("RTSP 虚拟设备已创建/存在", "id", rtspDevice.ID)
	}

	// 迁移 RTMP 推流数据
	if hasStreamPushs && rtmpDevice != nil {
		var pushes []StreamPush
		if err := db.WithContext(ctx).Find(&pushes).Error; err != nil {
			slog.Error("查询 stream_pushs 失败", "err", err)
			return err
		}

		migratedCount := 0
		for _, p := range pushes {
			// 检查是否已存在相同的通道
			var existing ipc.Channel
			if err := db.WithContext(ctx).Where("app = ? AND stream = ?", p.App, p.Stream).First(&existing).Error; err == nil {
				slog.Debug("通道已存在，跳过", "app", p.App, "stream", p.Stream)
				continue
			}

			channelID := uni.UniqueID(bz.IDPrefixRTMP)
			channel := ipc.Channel{
				ID:        channelID,
				DID:       rtmpDevice.ID,
				DeviceID:  rtmpDevice.ID, // RTMP 类型使用通道 ID 作为 device_id（唯一索引）
				ChannelID: channelID,     // RTMP 类型使用通道 ID 作为 channel_id
				Name:      p.Name,
				Type:      ipc.TypeRTMP,
				App:       p.App,
				Stream:    p.Stream,
				IsOnline:  p.Status == "PUSHING", // 根据旧状态设置在线状态
				Config: ipc.StreamConfig{
					IsAuthDisabled: p.IsAuthDisabled,
					Session:        p.Session,
					PushedAt:       p.PushedAt,
					StoppedAt:      p.StoppedAt,
					MediaServerID:  p.MediaServerID,
				},
				CreatedAt: p.CreatedAt,
				UpdatedAt: orm.Now(),
			}

			if err := db.WithContext(ctx).Create(&channel).Error; err != nil {
				slog.Error("迁移 RTMP 通道失败", "err", err, "stream", p.Stream)
				continue
			}
			migratedCount++
		}
		slog.Info("RTMP 数据迁移完成", "total", len(pushes), "migrated", migratedCount)
	}

	// 迁移 RTSP 代理数据
	if hasStreamProxys && rtspDevice != nil {
		var proxies []StreamProxy
		if err := db.WithContext(ctx).Find(&proxies).Error; err != nil {
			slog.Error("查询 stream_proxys 失败", "err", err)
			return err
		}

		migratedCount := 0
		for _, p := range proxies {
			// 检查是否已存在相同的通道
			var existing ipc.Channel
			if err := db.WithContext(ctx).Where("app = ? AND stream = ?", p.App, p.Stream).First(&existing).Error; err == nil {
				slog.Debug("通道已存在，跳过", "app", p.App, "stream", p.Stream)
				continue
			}

			channelID := uni.UniqueID(bz.IDPrefixRTSP)
			channel := ipc.Channel{
				ID:        channelID,
				DID:       rtspDevice.ID,
				DeviceID:  rtspDevice.ID, // RTSP 类型使用 ID 作为 device_id
				ChannelID: channelID,     // RTSP 类型使用 ID 作为 channel_id
				Name:      p.Name,
				Type:      ipc.TypeRTSP,
				App:       p.App,
				Stream:    p.Stream,
				IsOnline:  p.Pulling, // 根据旧拉流状态设置在线状态
				Config: ipc.StreamConfig{
					SourceURL:                 p.SourceURL,
					Transport:                 p.Transport,
					TimeoutS:                  p.TimeoutS,
					EnabledAudio:              p.EnabledAudio,
					EnabledRemoveNoneReader:   p.EnabledRemoveNoneReader,
					EnabledDisabledNoneReader: p.EnabledDisabledNoneReader,
					StreamKey:                 p.StreamKey,
					Enabled:                   p.Enabled,
				},
				CreatedAt: p.CreatedAt,
				UpdatedAt: orm.Now(),
			}

			if err := db.WithContext(ctx).Create(&channel).Error; err != nil {
				slog.Error("迁移 RTSP 通道失败", "err", err, "stream", p.Stream)
				continue
			}
			migratedCount++
		}
		slog.Info("RTSP 数据迁移完成", "total", len(proxies), "migrated", migratedCount)
	}

	slog.Info("数据迁移全部完成！旧表数据已保留，请手动确认后删除。")
	return nil
}
